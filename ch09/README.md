Here are some `curl` commands that match each of the routes in the teaching example web service:

---

**1. Root route** – Welcome message

```bash
curl http://localhost:8080/
```

**2. Time route** – Current server time in JSON

```bash
curl http://localhost:8080/time
```

**3. Echo with query parameter** – Echo a message back

```bash
curl "http://localhost:8080/echo?msg=HelloWorld"
```

**4. Echo with POST body** – Send JSON and get it echoed back

```bash
curl -X POST http://localhost:8080/echo \
     -H "Content-Type: application/json" \
     -d '{"name":"Mihalis","lang":"Go"}'
```

**5. Headers route** – See what headers your client sends

```bash
curl -H "X-Demo: true" http://localhost:8080/headers
```

**6. Health check route** – Service availability

```bash
curl http://localhost:8080/health
```

---


