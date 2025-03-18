import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './AllAppointments.css'; // Импортируйте стили для компонента

const AllAppointments = () => {
    const [appointments, setAppointments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchAppointments = async () => {
            try {
                const response = await axios.get('/admin/allAppointments');
                console.log('Fetched data:', response.data);
                setAppointments(response.data);
                setLoading(false);
            } catch (error) {
                console.error('Error fetching appointments:', error);
                setError('Failed to fetch appointments');
                setLoading(false);
            }
        };

        fetchAppointments();
    }, []);

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleDateString();
    };

    const formatTime = (timeString) => {
        const time = new Date(timeString);
        // Format to HH:mm (hours and minutes only)
        return time.toTimeString().split(' ')[0].substring(0, 5);
    };

    if (loading) return <div className="loading">Загрузка...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="appointments-container">
            <h3>Все Записи</h3>
            <table className="appointments-table">
                <thead>
                    <tr>
                        <th>ID Записи</th>
                        <th>Имя Пользователя</th>
                        <th>Контакт Пользователя</th>
                        <th>Дата Записи</th>
                        <th>Время Начала</th>
                        <th>Время Завершения</th>
                        <th>ID Услуги</th>
                        <th>ID Сотрудника</th>
                        <th>Статус</th>
                    </tr>
                </thead>
                <tbody>
                    {appointments.length > 0 ? (
                        appointments.map(appointment => (
                            <tr key={appointment.AppointmentID}>
                                <td>{appointment.AppointmentID}</td>
                                <td>{appointment.UserName}</td>
                                <td>{appointment.UserContact}</td>
                                <td>{formatDate(appointment.AppointmentDate)}</td>
                                <td>{formatTime(appointment.TimeStart)}</td>
                                <td>{formatTime(appointment.TimeEnd)}</td>
                                <td>{appointment.ServiceID}</td>
                                <td>{appointment.EmployeeID}</td>
                                <td>{appointment.Status ? 'Подтверждено' : 'Ожидание'}</td>
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan="9">Записи не найдены</td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    );
};

export default AllAppointments;
