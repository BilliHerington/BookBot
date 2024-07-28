import React, { useState } from 'react';
import axios from '../../../axios';
import './UploadNewSchedule.css'; // Импортируйте стили для компонента

const UploadNewSchedule = () => {
    const [file, setFile] = useState(null);

    const handleFileChange = (e) => {
        setFile(e.target.files[0]);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (file) {
            const formData = new FormData();
            formData.append('schedule', file); // Ensure the field name matches the backend

            try {
                await axios.post('/admin/uploadNewSchedule', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                });
                alert('Schedule uploaded successfully');
            } catch (error) {
                console.error('Error uploading schedule:', error);
                alert('Failed to upload schedule');
            }
        } else {
            alert('Please select a file first');
        }
    };

    return (
        <div className="upload-schedule-container">
            <h3>Загрузить Новое Расписание</h3>
            <form onSubmit={handleSubmit}>
                <label>
                    Выберите файл:
                    <input type="file" onChange={handleFileChange} required />
                </label>
                <button type="submit">Загрузить</button>
            </form>
        </div>
    );
};

export default UploadNewSchedule;
