import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './DeleteSchedule.css'; // Импортируйте стили для компонента

const DeleteSchedule = () => {
    const [schedules, setSchedules] = useState([]);
    const [selectedScheduleIds, setSelectedScheduleIds] = useState([]);
    const [selectAll, setSelectAll] = useState(false);

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

    const handleCheckboxChange = (scheduleID) => {
        setSelectedScheduleIds(prevSelectedIds => {
            if (prevSelectedIds.includes(scheduleID)) {
                return prevSelectedIds.filter(id => id !== scheduleID);
            } else {
                return [...prevSelectedIds, scheduleID];
            }
        });
    };

    const handleSelectAllChange = () => {
        if (selectAll) {
            setSelectedScheduleIds([]);
        } else {
            setSelectedScheduleIds(schedules.map(schedule => schedule.ScheduleID));
        }
        setSelectAll(!selectAll);
    };

    const handleDelete = async (e) => {
        e.preventDefault();

        if (selectedScheduleIds.length > 0) {
            try {
                await Promise.all(selectedScheduleIds.map(id => 
                    axios.delete(`/admin/deleteSchedule/${id}`)
                ));

                alert('Расписания удалены успешно');
                setSchedules(schedules.filter(schedule => !selectedScheduleIds.includes(schedule.ScheduleID)));
                setSelectedScheduleIds([]);
                setSelectAll(false);
            } catch (error) {
                console.error('Error deleting schedules:', error);
            }
        } else {
            alert('Пожалуйста, выберите хотя бы одно расписание для удаления');
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
        <div className="delete-schedules-container">
            <h3>Удаление Расписаний</h3>
            <form onSubmit={handleDelete}>
                <div className="schedules-table-container">
                    <table className="schedules-table">
                        <thead>
                            <tr>
                                <th>
                                    <input 
                                        type="checkbox" 
                                        checked={selectAll}
                                        onChange={handleSelectAllChange}
                                    />
                                    Выбрать все
                                </th>
                                <th>ID Расписания</th>
                                <th>ID Сотрудника</th>
                                <th>Дата Работы</th>
                                <th>Начало</th>
                                <th>Окончание</th>
                            </tr>
                        </thead>
                        <tbody>
                            {schedules.length > 0 ? (
                                schedules.map((schedule) => (
                                    <tr key={schedule.ScheduleID}>
                                        <td>
                                            <input 
                                                type="checkbox" 
                                                checked={selectedScheduleIds.includes(schedule.ScheduleID)}
                                                onChange={() => handleCheckboxChange(schedule.ScheduleID)}
                                            />
                                        </td>
                                        <td>{schedule.ScheduleID}</td>
                                        <td>{schedule.EmployeeID}</td>
                                        <td>{new Date(schedule.WorkDate).toLocaleDateString()}</td>
                                        <td>{formatTime(schedule.StartTime)}</td>
                                        <td>{formatTime(schedule.EndTime)}</td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan="6">Расписания отсутствуют</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
                <button type="submit">Удалить Выбранные Расписания</button>
            </form>
        </div>
    );
};

export default DeleteSchedule;
