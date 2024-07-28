import React, { useState, useEffect } from 'react';
import axios from '../../../axios';
import './AllServices.css'; // Import the CSS file

const AllServices = () => {
  const [services, setServices] = useState([]);

  useEffect(() => {
    const fetchServices = async () => {
      try {
        const response = await axios.get('/admin/allServices');
        setServices(response.data);
      } catch (error) {
        console.error('Error fetching services', error);
      }
    };

    fetchServices();
  }, []);

  return (
    <div className="services-container">
      <h2>All Services</h2>
      <table className="services-table">
        <thead>
          <tr>
            <th>Service ID</th>
            <th>Name</th>
            <th>Duration</th>
            <th>Default Price</th>
            <th>Pro Price</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {services.map(service => (
            <tr key={service.ServiceID}>
              <td>{service.ServiceID}</td>
              <td>{service.Name}</td>
              <td>{service.Duration}</td>
              <td>{service.DefaultPrice}</td>
              <td>{service.ProPrice}</td>
              <td>{service.Description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AllServices;
