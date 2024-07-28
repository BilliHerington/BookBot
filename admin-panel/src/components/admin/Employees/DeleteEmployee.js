import React, { useEffect, useState } from 'react';
import axios from '../../../axios'; // Убедитесь, что путь правильный
import './DeleteEmployee.css'; // Импортируем стили для компонента

const DeleteEmployee = () => {
    const [employees, setEmployees] = useState([]);
    const [selectedEmployeeIds, setSelectedEmployeeIds] = useState([]);

    // Функция для загрузки сотрудников
    useEffect(() => {
        const fetchEmployees = async () => {
            try {
                const response = await axios.get('/admin/allEmployees'); // Убедитесь, что правильный путь для получения сотрудников
                setEmployees(response.data);
            } catch (error) {
                console.error('Ошибка при получении сотрудников:', error);
            }
        };
        fetchEmployees();
    }, []);

    // Обработка изменения выбора чекбоксов
    const handleCheckboxChange = (employeeID) => {
        setSelectedEmployeeIds(prevSelectedIds => {
            if (prevSelectedIds.includes(employeeID)) {
                return prevSelectedIds.filter(id => id !== employeeID);
            } else {
                return [...prevSelectedIds, employeeID];
            }
        });
    };

    // Обработка удаления сотрудников
    const handleDelete = async (e) => {
        e.preventDefault();

        if (selectedEmployeeIds.length > 0) {
            try {
                await Promise.all(selectedEmployeeIds.map(id => 
                    axios.delete(`/admin/deleteEmployee/${id}`)
                ));

                alert('Сотрудники успешно удалены');
                setEmployees(employees.filter(employee => !selectedEmployeeIds.includes(employee.EmployeeID)));
                setSelectedEmployeeIds([]);
            } catch (error) {
                console.error('Ошибка при удалении сотрудников:', error);
            }
        } else {
            alert('Пожалуйста, выберите хотя бы одного сотрудника для удаления');
        }
    };

    return (
        <div className="delete-employee-container">
            <h2>Удаление сотрудников</h2>
            <form onSubmit={handleDelete} className="delete-employee-form">
                <label>Выберите сотрудников:</label>
                <div className="employee-list">
                    {employees.length > 0 ? (
                        employees.map((employee) => (
                            <div key={employee.EmployeeID} className="employee-item">
                                <input 
                                    type="checkbox" 
                                    checked={selectedEmployeeIds.includes(employee.EmployeeID)}
                                    onChange={() => handleCheckboxChange(employee.EmployeeID)}
                                />
                                <label>
                                    {employee.Name || 'Без имени'}
                                </label>
                            </div>
                        ))
                    ) : (
                        <div>Нет доступных сотрудников</div>
                    )}
                </div>
                <button type="submit" className="submit-button">Удалить выбранных сотрудников</button>
            </form>
        </div>
    );
};

export default DeleteEmployee;
