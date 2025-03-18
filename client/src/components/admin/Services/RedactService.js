import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './RedactService.css'; // Импортируем стили

const RedactService = () => {
    const [services, setServices] = useState([]);
    const [selectedService, setSelectedService] = useState(null);
    const [name, setName] = useState('');
    const [duration, setDuration] = useState('');
    const [def_price, setDefaultPrice] = useState('');
    const [pro_price, setProPrice] = useState('');
    const [description, setDescription] = useState('');

    useEffect(() => {
        const fetchServices = async () => {
            try {
                const response = await axios.get('/admin/allServices');
                setServices(response.data);
            } catch (error) {
                console.error('Error fetching services:', error);
            }
        };
        fetchServices();
    }, []);

    useEffect(() => {
        if (selectedService) {
            setName(selectedService.Name);
            setDuration(selectedService.Duration);
            setDefaultPrice(selectedService.DefaultPrice);
            setProPrice(selectedService.ProPrice);
            setDescription(selectedService.Description);
        } else {
            setName('');
            setDuration('');
            setDefaultPrice('');
            setProPrice('');
            setDescription('');
        }
    }, [selectedService]);

    const handleServiceChange = (e) => {
        const serviceId = parseInt(e.target.value, 10);
        const selected = services.find(service => service.ServiceID === serviceId);
        setSelectedService(selected || null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (selectedService) {
            try {
                const response = await axios.post('/admin/redactService', {
                    ID: selectedService.ServiceID,
                    Name: name,
                    Duration: duration,
                    DefPrice: parseFloat(def_price),
                    ProPrice: parseFloat(pro_price),
                    Description: description
                }, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (response.status === 200) {
                    alert('Service updated successfully');
                    const fetchServices = async () => {
                        try {
                            const response = await axios.get('/admin/allServices');
                            setServices(response.data);
                        } catch (error) {
                            console.error('Error fetching services:', error);
                        }
                    };
                    fetchServices();
                    setSelectedService(null);
                    setName('');
                    setDuration('');
                    setDefaultPrice('');
                    setProPrice('');
                    setDescription('');
                }
            } catch (error) {
                console.error('Error updating service:', error);
            }
        }
    };

    return (
        <div className="redact-service-container">
            <h2>Редактировать Услугу</h2>
            <form onSubmit={handleSubmit} className="redact-service-form">
                <label>
                    Выберите Услугу:
                    <select onChange={handleServiceChange} value={selectedService?.ServiceID || ''}>
                        <option value="">--Выберите услугу--</option>
                        {services.map((service) => (
                            <option key={service.ServiceID} value={service.ServiceID}>
                                {service.Name}
                            </option>
                        ))}
                    </select>
                </label>

                {selectedService && (
                    <>
                        <label>
                            Название:
                            <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
                        </label>
                        <label>
                            Продолжительность:
                            <input type="text" value={duration} onChange={(e) => setDuration(e.target.value)} required />
                        </label>
                        <label>
                            Обычная Цена:
                            <input type="number" step="0.01" value={def_price} onChange={(e) => setDefaultPrice(e.target.value)} required />
                        </label>
                        <label>
                            Премиум Цена:
                            <input type="number" step="0.01" value={pro_price} onChange={(e) => setProPrice(e.target.value)} required />
                        </label>
                        <label>
                            Описание:
                            <textarea value={description} onChange={(e) => setDescription(e.target.value)} required />
                        </label>
                        <button type="submit">Обновить Услугу</button>
                    </>
                )}
            </form>
        </div>
    );
};

export default RedactService;
