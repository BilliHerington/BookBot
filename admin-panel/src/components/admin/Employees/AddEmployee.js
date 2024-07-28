import React, { useState } from 'react';
import axios from '../../../axios';
import './AddEmployee.css'; // Импортируем стили для компонента

const AddEmployee = () => {
    const [name, setName] = useState('');
    const [level, setLevel] = useState('');
    const [contactNumber, setContactNumber] = useState('');
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('/admin/addEmployee', {
                name,
                level,
                contact_number: contactNumber // Отправляем в правильном формате
            });
            if (response.status === 200) {
                setSuccess('Employee added successfully');
                // Reset form fields
                setName('');
                setLevel('');
                setContactNumber('');
            }
        } catch (error) {
            console.error('Error adding employee:', error);
            setError('Error adding employee: ' + (error.response?.data?.error || error.message));
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
                    <input
                        id="level"
                        type="text"
                        value={level}
                        onChange={(e) => setLevel(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="contactNumber">Контактный номер:</label>
                    <input
                        id="contactNumber"
                        type="text" // Используйте type="text" чтобы избежать автоматического преобразования
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
