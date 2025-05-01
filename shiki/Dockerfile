# Use an official Node.js runtime as a parent image (Alpine for smaller size)
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /shiki-app

# Copy package.json and package-lock.json (or yarn.lock)
# This step leverages Docker layer caching
COPY package*.json ./

# Install pm2 globally and then production dependencies
RUN npm install -g pm2 && npm install --production

# Copy the rest of the application code
COPY . .

# Make port 3000 available within the Docker network (not exposed to host by default)
EXPOSE 3000

# Define the command to run the application using pm2
# -i max: Start instances based on available vCPUs
# --max-memory-restart: Restart instance if it exceeds this memory limit
CMD ["pm2-runtime", "server.js", "-i", "max", "--max-memory-restart", "512M"]