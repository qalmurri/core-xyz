import json
import random

def main():
    total_nodes = 2000
    map_nodes = []
    node_details = []

    base_prefixes = ["Sektor", "Pos Jaga", "Benteng", "Stasiun", "Markas"]
    mine_prefixes = ["Tambang", "Pengeboran", "Ekstraktor", "Situs"]
    suffixes = ["Alpha", "Beta", "Gamma", "Omega", "Prime", "Nova", "Zeta", "Nexus"]
    owners = ["Komandan_X", "Aliansi_Mawar", "Player_Satu", "Federasi", "Bajak Laut", "NPC_Trader"]

    print(f"[GENERATOR] Memulai pembuatan {total_nodes} node...")

    for node_id in range(1, total_nodes + 1):
        x = random.randint(-100, 100)
        y = random.randint(-100, 100)
        z = random.randint(-100, 100)
        level = random.randint(1, 4)
        type_node = random.randint(0, 1) 
        
        patrols = 0
        if type_node == 0:
            patrols = random.randint(0, 3)
        else:
            if random.random() > 0.7:
                patrols = 1
                
        # Simpan ke memori untuk JSON Peta
        map_nodes.append([node_id, x, y, z, level, type_node, patrols])
        
        suffix = random.choice(suffixes)
        owner = random.choice(owners)
        points = (random.randint(0, 49) + 10) * level
        
        if type_node == 0:
            name = f"{random.choice(base_prefixes)} {suffix}-{node_id}"
            type_name = "Pangkalan Utama"
            desc = f"Markas pertahanan level {level} yang dikendalikan oleh faksi kuat. Memiliki {patrols} armada yang bersiaga."
        else:
            name = f"{random.choice(mine_prefixes)} {suffix}-{node_id}"
            type_name = "Tambang Sumber Daya"
            desc = f"Fasilitas ekstraksi mineral level {level} yang beroperasi secara otomatis. Sangat berharga untuk mengumpulkan material."
            
        # Simpan ke memori untuk JSON Detail
        node_details.append({
            "id": node_id,
            "name": name,
            "owner": owner,
            "points": points,
            "type_name": type_name,
            "description": desc
        })
        
    with open("map_nodes.json", "w") as map_file:
        json.dump(map_nodes, map_file, indent=2)
        
    with open("node_details.json", "w") as detail_file:
        json.dump(node_details, detail_file, indent=2)

    # 5. Generate commands.txt untuk simulasi WebSockets Traffic
    print("[GENERATOR] Membuat skenario lalu lintas ke commands.txt...")
    commands_list = ["// Buka Console Browser (F12), lalu copy-paste baris di bawah ini untuk memicu simulasi lalu lintas!\n"]
    
    # Kita buat 50 rute acak antar planet yang baru saja digenerate
    for _ in range(1000):
        node_a = random.choice(map_nodes)
        node_b = random.choice(map_nodes)
        
        # Pastikan tidak memilih planet yang sama sebagai asal dan tujuan
        while node_a[0] == node_b[0]:
            node_b = random.choice(map_nodes)
            
        xa, ya, za = node_a[1], node_a[2], node_a[3]
        xb, yb, zb = node_b[1], node_b[2], node_b[3]
        
        duration = random.randint(10, 20) # Durasi perjalanan acak 3 - 10 detik
        
        # Format string JS: triggerRoute({x: 1, y: 2, z: 3}, {x: 4, y: 5, z: 6}, 5);
        cmd = f"triggerRoute({{x: {xa}, y: {ya}, z: {za}}}, {{x: {xb}, y: {yb}, z: {zb}}}, {duration});"
        commands_list.append(cmd)
        
    with open("commands.txt", "w") as cmd_file:
        cmd_file.write("\n".join(commands_list))

    print("[GENERATOR] Selesai! map_nodes.json, node_details.json, dan commands.txt berhasil dibuat.")

if __name__ == "__main__":
    main()
