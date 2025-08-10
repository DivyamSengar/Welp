"""
This script is used for parsing the durations of individual RPC requests from
kubernetes log files. 

The intended usage is to pipe kubernetes log files as input, and direct it towards 
this script. The output is a CSV file that can be put into Excel / Pandas / Google Sheets
for analysis.

Example:
    $ kubectl get pods
    NAME                                     READY   STATUS    RESTARTS   AGE
    detail-fdbcdf5bc-bb5hs                   1/1     Running   0          1m
    frontend-779c6bd97f-f4844                1/1     Running   0          1m
    ...
    reservation-7954cc84fd-958mf             1/1     Running   0          1m
    reviews-69997d7f97-cx2p9                 1/1     Running   0          1m
    $ kubectl logs detail-fdbcdf5bc-bb5hs | python parse_logs.py detail-fdbcdf5bc-bb5hs 
    Parsed data has been written to detail-fdbcdf5bc-bb5hs-logs.csv
"""
import re
import sys
import csv
from datetime import datetime

if len(sys.argv) < 2:
    print("Usage: python parse_logs.py <log_file.csv>")
    sys.exit(1)

pod_name = sys.argv[1]
csv_filename = f"{pod_name}-logs.csv"
pattern = re.compile(r'^(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (\S+);(.+);(.+);(.+);(.+);(.+)$')

header = ["date", "time", "rpc", "input", "output", "error", "start", "duration(µs)"]
# print(sys.argv)
# Open the CSV file for writing
with open(csv_filename, mode='w', newline='') as csv_file:
    csv_writer = csv.writer(csv_file)
    csv_writer.writerow(header)  # Write the header to the CSV file

    try:
        for line in sys.stdin:
            match = pattern.match(line)
            if match:
                datetime_tokens = datetime.strptime(match.group(1), "%Y/%m/%d %H:%M:%S")
                rpc_token = "/".join(match.group(2).split('/')[-2:])  # Extract last two parts of the rpc token
                input_token = match.group(3)
                output_token = match.group(4)
                error_token = match.group(5)
                start_token = match.group(6)
                duration_token = match.group(7)

                date = datetime_tokens.strftime("%m/%d/%Y")
                time = datetime_tokens.strftime("%H:%M:%S")

                # Write the parsed tokens to the CSV file
                csv_writer.writerow([date, time, rpc_token, input_token, output_token, error_token, start_token, duration_token])
    except KeyboardInterrupt:
        pass

print(f"Parsed data has been written to {csv_filename}")
