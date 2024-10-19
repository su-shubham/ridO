# Dockerfile.frontend

# Step 1: Build the frontend application
FROM node:20 as build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Build the application
RUN npm run build


# Expose the port on which the app runs
EXPOSE 80

