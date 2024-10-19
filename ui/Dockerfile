# get the base node image
FROM node:alpine as builder

# set the working dir for container
WORKDIR /frontend

# copy the json file first
COPY ./package.json /frontend
COPY ./package-lock.json /frontend

# extend npm file sizes
RUN export NODE_OPTIONS=--max_old_space_size=4096

# install npm dependencies
RUN npm cache clean --force
RUN npm ci --legacy-peer-deps

# copy other project files
COPY . .

# build the folder
#ARG REACT_APP_ENVIRONMENT
#ARG REACT_APP_BACKEND_URL
#RUN REACT_APP_ENVIRONMENT="$REACT_APP_ENVIRONMENT" REACT_APP_BACKEND_URL="$REACT_APP_BACKEND_URL" npm run build
RUN npm run build

# Handle Nginx
FROM nginx
COPY --from=builder /frontend/build /usr/share/nginx/html
COPY ./docker/nginx/default.conf /etc/nginx/conf.d/default.conf
