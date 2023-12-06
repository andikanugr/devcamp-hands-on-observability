# Devcamp-2023 Observability Hands-on

## Instructions
1. Simple instrumentation 
   - Install prometheus client instrumentation for client 
   - Add prometheus client instrumentation for server 
   - Try to run the server then open http://localhost:8080/metrics
2. Instrumentation with labels
   - Add labels to the instrumentation
   - Try to run the server then open http://localhost:8080/metrics
3. Prometheus dashboard
   - Create prometheus config file
   - Add prometheus to docker-compose
   - Try to run the server then open http://localhost:9090
   - Try to run the server then open http://localhost:9090/targets
   - Try to execute some queries
4. Try node exporter
   - Add node exporter to docker-compose
   - Try to run the server then open http://localhost:9090/targets
   - Try to execute some queries
5. Try postgres exporter
   - Add postgres exporter to docker-compose
   - Try to run the server then open http://localhost:9090/targets
   - Try to execute some queries