docker run -d \
--restart unless-stopped \
--name=mailpit \
-p 8025:8025 \
-p 1025:1025 \
axllent/mailpit

Failed deliveries are written to `dead-letter.csv` after the worker finishes
processing the recipient. The file includes the recipient, failure reason
(`template` or `smtp`), error details, and UTC timestamp.