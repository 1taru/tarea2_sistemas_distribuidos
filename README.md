# Gestion de notificaciones mediante KAFKA
Este proyecto para el montaje del docker compose junto a kafka usa la imagen de zookeper confluentinc/cp-zookeeper:latest y la imagen de kafka confluentinc/cp-kafka:latest en donde es necesario para iniciar el sistema de apache kafka
## Preparando el servicio
* primero que todo es necesario conectarse por terminal hacia la carpeta descargada llamada KAFKA-PROCESSING-SYSTEM en donde es necesario abrir 5 terminales
* en la primera terminal hay que usar para el levantamiento del servidor de kafka:
  ```bash
  docker-compose up --build --remove-orphans
  ```
* Luego en las otras 3 terminales, hay que hacer:
  para la segunda terminal:
  ```bash
  cd processor
  ```
  para la tercera terminal:
  ```bash
  cd producer
  ```
  para la cuarta terminal:
  ```bash
  cd consumer
  ```
  para la quinta terminal:
  ```bash
  cd client
  ```
* una vez hecho esto, hay que usar go run main.go en cada una de 4 ultimas terminales (tiene que ser la ultima la de la carpeta CLIENTE)
## Estudiantes
* Carlos Ruiz
* Nicolas Fernandez
