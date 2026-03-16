import time
import os
import psutil
from fastapi import FastAPI

app = FastAPI()

app.state.counter=0

@app.get("/{x}mb-{y}s")
def read_root(x, y):
    start = time.perf_counter()
    app.state.counter += 1
    x = int(x)
    y = int(y)
    # 100_000_000_000 = 100 gb
    # 100 mb
    occupated_space_mb = x * 1_000_000 # 1 mb
    time_out = y # in seconds
    strl = ' ' * occupated_space_mb
    time.sleep(time_out)
    app.state.counter-=1
    return {"timeout": time_out, "occupated_space": occupated_space_mb, "time": time.perf_counter()-start}

@app.get("/{x}fib")
def calculate_fib(x):
    start = time.perf_counter()
    app.state.counter += 1
    x = int(x)

    def fibrec(n):
        if n <= 2:
            return 1
        return fibrec(n-1) + fibrec(n-2)
    
    fibres = fibrec(x)
    app.state.counter -= 1
    return {"fib": fibres, "time": time.perf_counter()-start}

@app.get("/healthcheck")
def healthcheck():
    current_pid = os.getpid()
    cur_process = psutil.Process(current_pid)
    return {
        "counter": app.state.counter,
        "memory_info": f"{(cur_process.memory_info().rss) / 1024**2} MB",
        "cpu_usage": f"{cur_process.cpu_percent(interval=1)}%"
    }