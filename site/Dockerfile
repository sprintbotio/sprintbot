# Pull nginx base image
FROM nginx:latest

# Expost port 80
EXPOSE 80

# Copy custom configuration file from the current directory
COPY nginx-config/production.conf /etc/nginx/sites-enabled/sprintbot.conf

# Copy static assets into var/www
COPY ./index.html /var/www

# Start up nginx server
CMD ["nginx"]