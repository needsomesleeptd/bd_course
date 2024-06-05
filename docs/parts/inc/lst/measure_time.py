import requests
import time
import pandas as pd
import threading
import random

headers = {
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6Ik\
    pXVCJ9.eyJJRCI6MywiUm9sZSI6MiwiZXhwcmlyZXMiOiIyMDI0LTA1L\
    TIwVDExOjE2OjU3LjkyMzE0ODk5NSswMzowMCIsImxvZ2luIjoiYWRtaW4if\
    Q.kLhEqSSF8yK8GkCjtG7OixtgpDg8dicmM3F9Jk680B0',
    'Content-Type': 'application/json',
    'Cookie': 'auth_jwt="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.\
    eyJJRCI6MywiUm9sZSI6MiwiZXhwcmlyZXMiOiIyMDI0LTA1LTIwVDExOjE2O\
    jU3LjkyMzE0ODk5NSswMzowMCIsImxvZ2luIjoiYWRtaW4ifQ.kLhEqSSF8yK8\
    GkCjtG7OixtgpDg8dicmM3F9Jk680B0"'
}

urlTest = 'http://localhost:8080/annot/get'

mutex = threading.Lock()

def send_requests(url, n, interval, results,attempt):
    response_times = []
    rand_id = random.randint(1000,14000)
    start_time = time.time()
    response = requests.post(url, headers=headers, json={"id": rand_id})
    end_time = time.time()
    if response.status_code == 200:
        response_json = response.json()
        extracted_id = response_json.get('id')
        print("Extracted ID:", extracted_id)
    else:
        print("extract of ID failed", extracted_id)
    mutex.acquire()
    response_time = end_time - start_time
    response_times.append([response_time,n,attempt])
    mutex.release()
       
   

    results.extend(response_times)  # Store the response times in the results list

# Main function to send requests n times with a specified interval
def send_requests_n_times(url, n_requests, n_times, interval_seconds):
    results = []
    threads = []
    for run_id in range(times_block_reqs_ran):
        for n_time in range(n_times):
            
            #send_requests(url, n_requests, interval_seconds, results,n_time)
            thread = threading.Thread(target=send_requests, args=(url, n_times, interval_seconds, results,run_id))
            threads.append(thread)

          
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()
        time.sleep(1)
        threads = []
    return results



if __name__ == "__main__":
    times_block_reqs_ran = 40  # Replace with the total number of requests to be sent
    res_df = pd.DataFrame()
   
    for  n_requests_per_second in range(10,300 + 1,10):
        if n_requests_per_second % 10 == 0:
            print(f"serving {n_requests_per_second} requests per secound")

        interval_seconds = 1 / n_requests_per_second

        response_times = send_requests_n_times(urlTest, times_block_reqs_ran, n_requests_per_second, interval_seconds)

        df = pd.DataFrame({'Response Time': response_times})
        df_exploeded = df.explode('Response Time')

        # Split the response times into 'Value' and 'Time' columns
        all_vals_list = df_exploeded['Response Time'].tolist()
        Times = [all_vals_list[i] for i in range(len(all_vals_list)) if i % 3 == 0]
        ReqsCount = [all_vals_list[i] for i in range(len(all_vals_list)) if i % 3 == 1]
        Attempt = [all_vals_list[i] for i in range(len(all_vals_list)) if i % 3 == 2]
        df["Times"] = Times
        df["ReqsCount"] = ReqsCount
        df["Attempt"] = Attempt
        df = df.drop('Response Time', axis=1)

        res_df = pd.concat([res_df, df])

# Save the modified DataFrame to a CSV file
res_df.to_csv("response_data.csv", index=False)