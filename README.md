# Devcamp-2023 Observability Hands-on

## Instructions
1. Formatted logs
   - Install logrus 
   - Set log to file
   - Make sure you have a log folder
2. Loki and Promtail
   - Install Loki and Promtail on docker compose
   - Create Promtail config to scrape logs from the container
   - Update Dockerfile to include a log volume
3. Connect to Grafana
   - Add loki as a datasource
   - Try to query logs
   - Create a dashboard to view logs
