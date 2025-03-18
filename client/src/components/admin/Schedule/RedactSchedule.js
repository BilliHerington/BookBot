import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './RedactSchedule.css'; // Импортируйте стили для компонента

const RedactSchedule = () => {
    const [schedules, setSchedules] = useState([]);
    const [selectedSchedule, setSelectedSchedule] = useState(null);
    const [employeeID, setEmployeeID] = useState('');
    const [workDate, setWorkDate] = useState('');
    const [startTime, setStartTime] = useState('');
    const [endTime, setEndTime] = useState('');

    useEffect(() => {
        const fetchSchedules = async () => {
            try {
                const response = await axios.get('/admin/allSchedules');
                const formattedSchedules = response.data.map(schedule => ({
                    ...schedule,
                    DisplayText: `${schedule.EmployeeID} - ${new Date(schedule.WorkDate).toLocaleDateString()}`
                }));
                setSchedules(formattedSchedules);
            } catch (error) {
                console.error('Error fetching schedules:', error.response ? error.response.data : error.message);
            }
        };
        fetchSchedules();
    }, []);

    useEffect(() => {
        if (selectedSchedule) {
            const date = new Date(selectedSchedule.WorkDate);
            const formattedDate = date.toISOString().split('T')[0]; // Format: yyyy-MM-dd
            const formattedStartTime = selectedSchedule.StartTime.split('T')[1].substring(0, 5); // Format: HH:mm
            const formattedEndTime = selectedSchedule.EndTime.split('T')[1].substring(0, 5); // Format: HH:mm

            setEmployeeID(selectedSchedule.EmployeeID.toString());
            setWorkDate(formattedDate);
            setStartTime(formattedStartTime);
            setEndTime(formattedEndTime);
        } else {
            setEmployeeID('');
            setWorkDate('');
            setStartTime('');
            setEndTime('');
        }
    }, [selectedSchedule]);

    const handleScheduleChange = (e) => {
        const scheduleId = parseInt(e.target.value, 10);
        const selected = schedules.find(schedule => schedule.ScheduleID === scheduleId);
        setSelectedSchedule(selected || null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (selectedSchedule) {
            try {
                // Ensure all fields are filled
                if (!employeeID || !workDate || !startTime || !endTime) {
                    alert("All fields must be filled.");
                    return;
                }

                const response = await axios.post('/admin/redactSchedule', {
                    id: selectedSchedule.ScheduleID,
                    employee_id: parseInt(employeeID, 10),
                    work_date: workDate,
                    time_start: `${startTime}:00`,
                    time_end: `${endTime}:00`
                }, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (response.status === 200) {
                    alert('Schedule updated successfully');
                    const fetchSchedules = async () => {
                        try {
                            const response = await axios.get('/admin/allSchedules');
                            const formattedSchedules = response.data.map(schedule => ({
                                ...schedule,
                                DisplayText: `${schedule.EmployeeID} - ${new Date(schedule.WorkDate).toLocaleDateString()}`
                            }));
                            setSchedules(formattedSchedules);
                        } catch (error) {
                            console.error('Error fetching schedules:', error.response ? error.response.data : error.message);
                        }
                    };
                    fetchSchedules();
                    setSelectedSchedule(null);
                    setEmployeeID('');
                    setWorkDate('');
                    setStartTime('');
                    setEndTime('');
                }
            } catch (error) {
                console.error('Error updating schedule:', error.response ? error.response.data : error.message);
            }
        }
    };

    return (
        <div className="redact-schedule-container">
            <h3>Редактировать Расписание</h3>
            <form onSubmit={handleSubmit}>
                <label>
                    Выберите расписание:
                    <select onChange={handleScheduleChange} value={selectedSchedule?.ScheduleID || ''}>
                        <option value="">--Выберите расписание--</option>
                        {schedules.map((schedule) => (
                            <option key={schedule.ScheduleID} value={schedule.ScheduleID}>
                                {schedule.DisplayText}
                            </option>
                        ))}
                    </select>
                </label>

                {selectedSchedule && (
                    <>
                        <label>
                            ID сотрудника:
                            <input type="text" value={employeeID} onChange={(e) => setEmployeeID(e.target.value)} required />
                        </label>
                        <label>
                            Дата работы:
                            <input type="date" value={workDate} onChange={(e) => setWorkDate(e.target.value)} required />
                        </label>
                        <label>
                            Время начала:
                            <input type="time" value={startTime} onChange={(e) => setStartTime(e.target.value)} required />
                        </label>
                        <label>
                            Время окончания:
                            <input type="time" value={endTime} onChange={(e) => setEndTime(e.target.value)} required />
                        </label>
                        <button type="submit">Обновить Расписание</button>
                    </>
                )}
            </form>
        </div>
    );
};

export default RedactSchedule;
