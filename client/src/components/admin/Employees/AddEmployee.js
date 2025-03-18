import React, { useState, useEffect } from 'react';
import axios from '../../../axios';
import './AddEmployee.css';

const AddEmployee = () => {
    const [name, setName] = useState('');
    const [level, setLevel] = useState('');
    const [contactNumber, setContactNumber] = useState('');
    const [levels, setLevels] = useState([]); // <-- Гарантируем, что это массив
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);

    useEffect(() => {
        const fetchLevels = async () => {
            try {
                const response = await axios.get('/admin/employeeLevels');
                setLevels(response.data.levels || []); // <-- Убедимся, что `levels` всегда массив
            } catch (error) {
                console.error('Ошибка при загрузке уровней:', error);
                setLevels([]); // <-- В случае ошибки не даем `undefined`
                setError('Не удалось загрузить уровни сотрудников');
            }
        };

        fetchLevels();
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('/admin/addEmployee', {
                name,
                level,
                contact_number: contactNumber
            });
            if (response.status === 200) {
                setSuccess('Сотрудник успешно добавлен');
                setName('');
                setLevel('');
                setContactNumber('');
            }
        } catch (error) {
            console.error('Ошибка добавления сотрудника:', error);
            setError('Ошибка добавления сотрудника: ' + (error.response?.data?.error || error.message));
        }
    };

    return (
        <div className="employee-container">
            <h2>Добавить Сотрудника</h2>
            <form onSubmit={handleSubmit} className="employee-form">
                <div className="form-group">
                    <label htmlFor="name">Имя:</label>
                    <input
                        id="name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="level">Уровень:</label>
                    <select id="level" value={level} onChange={(e) => setLevel(e.target.value)} required>
                        <option value="">Выберите уровень</option>
                        {levels.map((lvl) => (
                            <option key={lvl} value={lvl}>{lvl}</option>
                        ))}
                    </select>
                </div>
                <div className="form-group">
                    <label htmlFor="contactNumber">Контактный номер:</label>
                    <input
                        id="contactNumber"
                        type="text"
                        value={contactNumber}
                        onChange={(e) => setContactNumber(e.target.value)}
                        required
                    />
                </div>
                <button type="submit" className="submit-button">Добавить Сотрудника</button>
            </form>
            {success && <p className="success-message">{success}</p>}
            {error && <p className="error-message">{error}</p>}
        </div>
    );
};

export default AddEmployee;
