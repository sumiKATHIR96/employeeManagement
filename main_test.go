package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCRUDOperations(t *testing.T) {
	cac := NewEmployeeCache()

	emp := Employee{
		ID:       "1",
		Name:     "sumithra",
		Position: "Engineer",
		Salary:   50000,
	}
	createdEmp, err := cac.CreateEmployee(emp)
	fmt.Println(createdEmp)
	assert.NoError(t, err)
	assert.Equal(t, emp, createdEmp)

	// Test GetEmployeeByID
	retrievedEmp, err := cac.GetEmployeeByID("1")
	assert.NoError(t, err)
	assert.Equal(t, emp, retrievedEmp)

	// Test UpdateEmployee
	updatedEmp, err := cac.UpdateEmployee("1", "sumithra", "Engineer", 70000)
	assert.NoError(t, err)
	assert.Equal(t, "sumithra", updatedEmp.Name)
	assert.Equal(t, "Engineer", updatedEmp.Position)
	assert.Equal(t, 70000, updatedEmp.Salary)

	// Test DeleteEmployee
	err = cac.DeleteEmployee("1")
	assert.NoError(t, err)

	// Verify employee is deleted
	_, err = cac.GetEmployeeByID("1")
	assert.Error(t, err)
}
