import React, { useEffect, useState } from 'react';
import axios from '../../../axios';
import './DeleteService.css'; // Импортируем CSS файл

const DeleteService = () => {
    const [services, setServices] = useState([]);
    const [selectedServiceIds, setSelectedServiceIds] = useState([]);

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

    const handleCheckboxChange = (serviceID) => {
        setSelectedServiceIds(prevSelectedIds => {
            if (prevSelectedIds.includes(serviceID)) {
                return prevSelectedIds.filter(id => id !== serviceID);
            } else {
                return [...prevSelectedIds, serviceID];
            }
        });
    };

    const handleDelete = async (e) => {
        e.preventDefault();

        if (selectedServiceIds.length > 0) {
            try {
                await Promise.all(selectedServiceIds.map(id => 
                    axios.delete(`/admin/deleteService/${id}`)
                ));

                alert('Services deleted successfully');
                setServices(services.filter(service => !selectedServiceIds.includes(service.ServiceID)));
                setSelectedServiceIds([]);
            } catch (error) {
                console.error('Error deleting services:', error);
            }
        } else {
            alert('Please select at least one service to delete');
        }
    };

    return (
        <div className="delete-service-container">
            <h2>Удаление Услуг</h2>
            <form onSubmit={handleDelete} className="delete-service-form">
                <label className="select-label">Выберите Услуги:</label>
                <div className="service-list">
                    {services.length > 0 ? (
                        services.map((service) => (
                            <div key={service.ServiceID} className="service-item">
                                <input 
                                    type="checkbox" 
                                    id={`service-${service.ServiceID}`}
                                    checked={selectedServiceIds.includes(service.ServiceID)}
                                    onChange={() => handleCheckboxChange(service.ServiceID)}
                                />
                                <label htmlFor={`service-${service.ServiceID}`} className="service-label">
                                    {service.Name || 'Без названия'}
                                </label>
                            </div>
                        ))
                    ) : (
                        <div className="no-services">Нет доступных услуг</div>
                    )}
                </div>
                <button type="submit" className="delete-button">Удалить Выбранные Услуги</button>
            </form>
        </div>
    );
};

export default DeleteService;
