import React, { useState } from 'react';
import axios from '../../../axios';
import './AddService.css'; // Import the CSS file

const AddService = () => {
    const [name, setName] = useState('');
    const [duration, setDuration] = useState('');
    const [defPrice, setDefPrice] = useState('');
    const [proPrice, setProPrice] = useState('');
    const [description, setDescription] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await axios.post('/admin/addService', {
                Name: name,
                Duration: duration,
                DefPrice: parseFloat(defPrice), // Ensure the price is sent as a float
                ProPrice: parseFloat(proPrice), // Ensure the price is sent as a float
                Description: description
            });
            alert('Service added successfully');
            // Reset form fields
            setName('');
            setDuration('');
            setDefPrice('');
            setProPrice('');
            setDescription('');
        } catch (error) {
            console.error('Error adding service:', error);
        }
    };

    return (
        <div className="service-container">
            <h3>Add Service</h3>
            <form onSubmit={handleSubmit}>
                <label>
                    Name:
                    <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
                </label>
                <label>
                    Duration:
                    <input type="text" value={duration} onChange={(e) => setDuration(e.target.value)} required />
                </label>
                <label>
                    Default Price:
                    <input type="number" step="0.01" value={defPrice} onChange={(e) => setDefPrice(e.target.value)} required />
                </label>
                <label>
                    Pro Price:
                    <input type="number" step="0.01" value={proPrice} onChange={(e) => setProPrice(e.target.value)} required />
                </label>
                <label>
                    Description:
                    <textarea value={description} onChange={(e) => setDescription(e.target.value)} required />
                </label>
                <button type="submit">Add Service</button>
            </form>
        </div>
    );
};

export default AddService;
