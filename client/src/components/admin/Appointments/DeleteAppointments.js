import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './DeleteAppointments.css'; // Импортируйте стили для компонента

const DeleteAppointments = () => {
    const [appointments, setAppointments] = useState([]);
    const [selectedAppointmentIds, setSelectedAppointmentIds] = useState([]);

    useEffect(() => {
        const fetchAppointments = async () => {
            try {
                const response = await axios.get('/admin/allAppointments');
                setAppointments(response.data);
            } catch (error) {
                console.error('Error fetching appointments:', error);
            }
        };
        fetchAppointments();
    }, []);

    const handleCheckboxChange = (appointmentID) => {
        setSelectedAppointmentIds(prevSelectedIds => {
            if (prevSelectedIds.includes(appointmentID)) {
                return prevSelectedIds.filter(id => id !== appointmentID);
            } else {
                return [...prevSelectedIds, appointmentID];
            }
        });
    };

    const handleDelete = async (e) => {
        e.preventDefault();

        if (selectedAppointmentIds.length > 0) {
            try {
                await Promise.all(selectedAppointmentIds.map(id => 
                    axios.delete(`/admin/deleteAppointment/${id}`)
                ));

                alert('Записи удалены успешно');
                setAppointments(appointments.filter(appointment => !selectedAppointmentIds.includes(appointment.AppointmentID)));
                setSelectedAppointmentIds([]);
            } catch (error) {
                console.error('Error deleting appointments:', error);
            }
        } else {
            alert('Пожалуйста, выберите хотя бы одну запись для удаления');
        }
    };

    const formatTime = (timeString) => {
        if (!timeString) {
            return 'Invalid Time';
        }

        try {
            const [timePart] = timeString.split('T')[1]?.split('Z');
            const [hours, minutes] = timePart.split(':');
            return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`;
        } catch (error) {
            return 'Invalid Time';
        }
    };

    return (
        <div className="delete-appointments-container">
            <h3>Удаление Записей</h3>
            <form onSubmit={handleDelete}>
                <table className="appointments-table">
                    <thead>
                        <tr>
                            <th>Выбрать</th>
                            <th>ID Записи</th>
                            <th>Имя Пользователя</th>
                            <th>ID Сотрудника</th>
                            <th>Дата Записи</th>
                            <th>Начало</th>
                            <th>Окончание</th>
                        </tr>
                    </thead>
                    <tbody>
                        {appointments.length > 0 ? (
                            appointments.map((appointment) => (
                                <tr key={appointment.AppointmentID}>
                                    <td>
                                        <input 
                                            type="checkbox" 
                                            checked={selectedAppointmentIds.includes(appointment.AppointmentID)}
                                            onChange={() => handleCheckboxChange(appointment.AppointmentID)}
                                        />
                                    </td>
                                    <td>{appointment.AppointmentID}</td>
                                    <td>{appointment.UserName || 'Не указано'}</td>
                                    <td>{appointment.EmployeeID}</td>
                                    <td>{new Date(appointment.AppointmentDate).toLocaleDateString()}</td>
                                    <td>{formatTime(appointment.TimeStart)}</td>
                                    <td>{formatTime(appointment.TimeEnd)}</td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan="7">Записи отсутствуют</td>
                            </tr>
                        )}
                    </tbody>
                </table>
                <button type="submit">Удалить Выбранные Записи</button>
            </form>
        </div>
    );
};

export default DeleteAppointments;
