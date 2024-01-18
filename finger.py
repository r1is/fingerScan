import json
import hashlib

# 指纹去重脚本

# 读取JSON文件
with open('data.json', 'r', encoding='utf-8') as file:
    data = json.load(file)

# 创建一个空集合用于存放hash值
hash_set = set()

# 创建一个空列表用于存放没有重复hash的项
unique_items = []

# 对列表中的每一项进行hash计算
for item in data["fingerprint"]:
    # 计算hash值
    item_hash = hashlib.sha256(json.dumps(item, sort_keys=True).encode('utf-8')).hexdigest()

    # 如果hash值不在集合中，则将其添加到集合和结果列表中
    if item_hash not in hash_set:
        hash_set.add(item_hash)
        unique_items.append(item)
# 将结果保存到新的JSON文件中
with open('unique_data.json', 'w', encoding='utf-8') as output_file:
    json.dump(unique_items, output_file, ensure_ascii=False)

# 打印结果列表
print(unique_items)