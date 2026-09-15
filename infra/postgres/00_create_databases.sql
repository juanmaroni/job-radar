-- This script is not for production. Replace the "changeMe".
CREATE USER tracker_service WITH PASSWORD 'changeMe';
CREATE DATABASE tracker_service OWNER tracker_service;

CREATE USER jobs_service WITH PASSWORD 'changeMe';
CREATE DATABASE jobs_service OWNER jobs_service;
