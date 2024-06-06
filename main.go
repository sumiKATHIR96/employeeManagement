package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"sync"
)

type Employee struct {
	ID       string `json:"id" `
	Name     string `json:"name" binding:"required"`
	Position string `json:"position" binding:"required"`
	Salary   int    `json:"salary" binding:"required"`
}

type EmployeeCache struct {
	sync.RWMutex
	employeeData map[string]Employee
}

func NewEmployeeCache() *EmployeeCache {
	return &EmployeeCache{
		employeeData: make(map[string]Employee),
	}
}

func (cache *EmployeeCache) CreateEmployee(emp Employee) (Employee, error) {
	cache.Lock()
	defer cache.Unlock()
	_, exists := cache.employeeData[emp.ID]
	if exists {
		return Employee{}, errors.New("employee already found")
	}
	cache.employeeData[emp.ID] = emp
	return emp, nil
}

func (cache *EmployeeCache) GetEmployeeByID(id string) (Employee, error) {
	cache.RLock()
	defer cache.RUnlock()
	emp, exists := cache.employeeData[id]
	if !exists {
		return Employee{}, errors.New("employee not found")
	}
	return emp, nil
}

func (cache *EmployeeCache) UpdateEmployee(id string, name, position string, salary int) (Employee, error) {
	cache.Lock()
	defer cache.Unlock()
	emp, exists := cache.employeeData[id]
	if !exists {
		return Employee{}, errors.New("employee not found")
	}
	emp.Name = name
	emp.Position = position
	emp.Salary = salary
	cache.employeeData[id] = emp
	return emp, nil
}

func (cache *EmployeeCache) DeleteEmployee(id string) error {
	cache.Lock()
	defer cache.Unlock()
	_, exists := cache.employeeData[id]
	if !exists {
		return errors.New("employee not found")
	}
	delete(cache.employeeData, id)
	return nil
}

var cache = NewEmployeeCache()

func ListEmployeesDetails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	employees := make([]Employee, 0, len(cache.employeeData))
	for _, emp := range cache.employeeData {
		employees = append(employees, emp)
	}

	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(employees) {
		start = len(employees)
	}
	if end >= len(employees) {
		end = len(employees)
	}

	paginatedEmployees := employees[start:end]
	c.JSON(http.StatusOK, paginatedEmployees)
}

func CreateEmployeeDetails(c *gin.Context) {
	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emp.ID = uuid.New().String()
	newEmp, err := cache.CreateEmployee(emp)
	if err != nil {
		c.JSON(http.StatusFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newEmp)
}

func GetEmployeeDetails(c *gin.Context) {
	id := c.Param("id")
	emp, err := cache.GetEmployeeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func UpdateEmployeeDetails(c *gin.Context) {
	id := c.Param("id")
	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEmp, err := cache.UpdateEmployee(id, emp.Name, emp.Position, emp.Salary)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedEmp)
}

func DeleteEmployeeDetails(c *gin.Context) {
	id := c.Param("id")
	err := cache.DeleteEmployee(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func main() {
	r := gin.Default()
	r.GET("/employee", ListEmployeesDetails)
	r.POST("/employee", CreateEmployeeDetails)
	r.GET("/employee/:id", GetEmployeeDetails)
	r.PUT("/employee/:id", UpdateEmployeeDetails)
	r.DELETE("/employee/:id", DeleteEmployeeDetails)

	r.Run(":8080")
}
