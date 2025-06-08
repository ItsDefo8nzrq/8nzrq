import os
import requests
from concurrent.futures import ThreadPoolExecutor, as_completed

def fetch_proxies(url):
    try:
        response = requests.get(url, timeout=10)
        # Assuming the proxies are listed line by line
        return response.text.splitlines()
    except Exception as e:
        print(f"Error fetching proxies from {url}: {e}")
        return []

def check_proxy(proxy):
    # Ensure the proxy URL has the proper format
    if not proxy.startswith("http://") and not proxy.startswith("https://"):
        proxy = "http://" + proxy
    try:
        response = requests.get('http://httpbin.org/ip', proxies={"http": proxy, "https": proxy}, timeout=5)
        if response.status_code == 200:
            return proxy
    except Exception:
        pass
    return None

def main():
    # Ensure that the "validproxies" file exists before starting
    script_dir = os.path.dirname(os.path.realpath(__file__))
    valid_proxies_path = os.path.join(script_dir, "validproxies")
    if not os.path.exists(valid_proxies_path):
        print("File 'validproxies' does not exist. Please create the file in the same folder as the proxy checker.")
        return

    # List of URLs with proxy lists
    urls = [
        'https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt',
        'https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-http.txt'
    ]
    
    all_proxies = set()
    for url in urls:
        proxies = fetch_proxies(url)
        all_proxies.update(proxies)
    
    print(f"Total proxies fetched: {len(all_proxies)}")
    
    valid_proxies = []
    # Use ThreadPoolExecutor for faster concurrent checking
    with ThreadPoolExecutor(max_workers=50) as executor:
        future_to_proxy = {executor.submit(check_proxy, proxy): proxy for proxy in all_proxies}
        for future in as_completed(future_to_proxy):
            result = future.result()
            if result:
                print(f'Valid Proxy: {result}')
                valid_proxies.append(result)
            else:
                # Uncomment the next line to see invalid proxies
                # print(f'Invalid Proxy: {future_to_proxy[future]}')
                pass

    # Write valid proxies to the existing file "validproxies"
    try:
        with open(valid_proxies_path, "w") as f:
            for proxy in valid_proxies:
                f.write(proxy + "\n")
        print(f"Total valid proxies: {len(valid_proxies)} written to {valid_proxies_path}")
    except Exception as e:
        print(f"Error writing valid proxies to file: {e}")

if __name__ == "__main__":
    main()