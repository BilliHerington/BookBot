import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './AllSchedules.css'; // Импортируйте стили для компонента

const AllSchedules = () => {
    const [schedules, setSchedules] = useState([]);

    useEffect(() => {
        const fetchSchedules = async () => {
            try {
                const response = await axios.get('/admin/allSchedules');
                setSchedules(response.data);
            } catch (error) {
                console.error('Error fetching schedules:', error);
            }
        };
        fetchSchedules();
    }, []);

    // Форматирование даты в формат dd/MM/yyyy
    const formatDate = (dateStr) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-GB'); // Форматирование даты
    };

    // Форматирование времени в формат HH:mm
    const formatTime = (timeStr) => {
        // Преобразуем строку времени в объект Date
        // Заменяем T и Z, чтобы корректно преобразовать строку в объект Date
        const date = new Date(timeStr.replace('Z', ''));
        // Проверяем, корректно ли преобразована строка времени
        if (isNaN(date.getTime())) {
            return 'Invalid Time'; // Возвращаем сообщение об ошибке, если преобразование не удалось
        }
        // Форматируем время в HH:mm
        return date.toTimeString().slice(0, 5);
    };

    return (
        <div className="all-schedules-container">
            <h3>Все Расписания</h3>
            {schedules.length > 0 ? (
                <div style={{ maxHeight: '500px', overflowY: 'auto' }}>
                    <table className="all-schedules-table">
                        <thead>
                            <tr>
                                <th>ID Расписания</th>
                                <th>ID Сотрудника</th>
                                <th>Дата Работы</th>
                                <th>Время Начала</th>
                                <th>Время Окончания</th>
                            </tr>
                        </thead>
                        <tbody>
                            {schedules.map((schedule) => (
                                <tr key={schedule.ScheduleID}>
                                    <td>{schedule.ScheduleID}</td>
                                    <td>{schedule.EmployeeID}</td>
                                    <td>{formatDate(schedule.WorkDate)}</td>
                                    <td>{formatTime(schedule.StartTime)}</td>
                                    <td>{formatTime(schedule.EndTime)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            ) : (
                <p>No schedules found.</p>
            )}
        </div>
    );
};

export default AllSchedules;
