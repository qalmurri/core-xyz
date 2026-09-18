import json
import random

colors = ["0x38bdf8", "0xfacc15", "0x2ecc71", "0xe879f9", "0xf97316"]
types = ["base", "tambang"]
nodes = []
used_coords = set()

for i in range(1, 501):
    while True:
        x = random.randint(-10, 10)
        y = random.randint(-10, 10)
        z = random.randint(-10, 10)
        if (x, y, z) not in used_coords:
            used_coords.add((x, y, z))
            break
    
    node_type = random.choice(types)
    level = random.randint(1, 4)
    
    patrols = []
    p_count = random.randint(0, 2)
    for p in range(p_count):
        patrols.append({
            "speed": round(0.8 + random.random() * 1.2, 2),
            "radius": round(0.35 + random.random() * 0.2, 2),
            "color": colors[p % len(colors)],
            "tilt": round(random.random() * 0.6, 2)
        })
        
    name = f"Tambang Kristal #{i}" if node_type == "tambang" else f"Desa Sektor #{i}"
    owner = "Netral" if node_type == "tambang" else f"Komandan_{random.randint(1, 100)}"
    
    nodes.append({
        "id": f"node-{i}",
        "name": name,
        "type": node_type,
        "level": level,
        "x": x,
        "y": y,
        "z": z,
        "patrols": patrols,
        "owner": owner,
        "points": level * 350
    })

with open("nodes.json", "w", encoding="utf-8") as f:
    json.dump(nodes, f, indent=2)

print("Berhasil membuat file nodes.json dengan 500 data!")