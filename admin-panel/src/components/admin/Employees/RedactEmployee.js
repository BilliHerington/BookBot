import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './RedactEmployee.css'; // Импортируем стили для компонента

const RedactEmployee = () => {
    const [employees, setEmployees] = useState([]);
    const [selectedEmployee, setSelectedEmployee] = useState(null);
    const [name, setName] = useState('');
    const [level, setLevel] = useState('');
    const [contactNumber, setContactNumber] = useState('');

    useEffect(() => {
        const fetchEmployees = async () => {
            try {
                const response = await axios.get('/admin/allEmployees');
                setEmployees(response.data);
            } catch (error) {
                console.error('Error fetching employees:', error);
            }
        };
        fetchEmployees();
    }, []);

    useEffect(() => {
        if (selectedEmployee) {
            setName(selectedEmployee.Name || '');
            setLevel(selectedEmployee.Level || '');
            setContactNumber(selectedEmployee.ContactNumber || '');
        } else {
            setName('');
            setLevel('');
            setContactNumber('');
        }
    }, [selectedEmployee]);

    const handleEmployeeChange = (e) => {
        const employeeId = parseInt(e.target.value, 10);
        const selected = employees.find(employee => employee.EmployeeID === employeeId);
        setSelectedEmployee(selected || null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (selectedEmployee) {
            try {
                const response = await axios.post('/admin/redactEmployee', {
                    id: selectedEmployee.EmployeeID,
                    name,
                    level,
                    contact_number: contactNumber
                }, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (response.status === 200) {
                    alert('Employee updated successfully');
                    const fetchEmployees = async () => {
                        try {
                            const response = await axios.get('/admin/allEmployees');
                            setEmployees(response.data);
                        } catch (error) {
                            console.error('Error refetching employees:', error);
                        }
                    };
                    fetchEmployees();
                    setSelectedEmployee(null);
                    setName('');
                    setLevel('');
                    setContactNumber('');
                }
            } catch (error) {
                console.error('Error updating employee:', error);
            }
        }
    };

    return (
        <div className="redact-employee-container">
            <h3>Редактировать Сотрудника</h3>
            <form onSubmit={handleSubmit}>
                <label>
                    Выберите сотрудника:
                    <select 
                        onChange={handleEmployeeChange} 
                        value={selectedEmployee?.EmployeeID || ''}
                    >
                        <option value="">--Выберите сотрудника--</option>
                        {employees.map((employee) => (
                            <option key={employee.EmployeeID} value={employee.EmployeeID}>
                                {employee.Name}
                            </option>
                        ))}
                    </select>
                </label>

                {selectedEmployee && (
                    <>
                        <label>
                            Имя:
                            <input 
                                type="text" 
                                value={name} 
                                onChange={(e) => setName(e.target.value)} 
                                required 
                            />
                        </label>
                        <label>
                            Уровень:
                            <input 
                                type="text" 
                                value={level} 
                                onChange={(e) => setLevel(e.target.value)} 
                                required 
                            />
                        </label>
                        <label>
                            Контактный номер:
                            <input 
                                type="text" 
                                value={contactNumber} 
                                onChange={(e) => setContactNumber(e.target.value)} 
                                required 
                            />
                        </label>
                        <button type="submit">Обновить Сотрудника</button>
                    </>
                )}
            </form>
        </div>
    );
};

export default RedactEmployee;
