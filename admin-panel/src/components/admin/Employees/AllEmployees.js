import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './AllEmployees.css'; // Импортируем стили для компонента

const AllEmployees = () => {
    const [employees, setEmployees] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchEmployees = async () => {
            try {
                const response = await axios.get('/admin/allEmployees');
                setEmployees(response.data);
                setLoading(false);
            } catch (error) {
                console.error('Ошибка при получении сотрудников:', error);
                setError(error);
                setLoading(false);
            }
        };

        fetchEmployees();
    }, []);

    if (loading) {
        return <div className="loading">Загрузка...</div>;
    }

    if (error) {
        return <div className="error">Ошибка: {error.message}</div>;
    }

    return (
        <div className="employees-container">
            <h2>Все сотрудники</h2>
            {employees.length === 0 ? (
                <p>Нет данных о сотрудниках</p>
            ) : (
                <table className="employees-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Имя</th>
                            <th>Уровень</th>
                            <th>Контактный номер</th>
                        </tr>
                    </thead>
                    <tbody>
                        {employees.map((employee) => (
                            <tr key={employee.EmployeeID}>
                                <td>{employee.EmployeeID}</td>
                                <td>{employee.Name}</td>
                                <td>{employee.Level}</td>
                                <td>{employee.ContactNumber}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default AllEmployees;
