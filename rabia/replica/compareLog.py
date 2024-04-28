import re

def compare_consensus_values(log1_content, log2_content):
    # 正規表現パターンを使用して、consensusValueの値を抽出する
    pattern = r'consensusValue:\s+(\d+)'
    
    # ログファイル1からconsensusValueの値を抽出する
    values1 = re.findall(pattern, log1_content)
    
    # ログファイル2からconsensusValueの値を抽出する
    values2 = re.findall(pattern, log2_content)
    
    # 抽出したconsensusValueの値が一致しているかどうかを確認する
    if values1 == values2:
        print("consensusValueの値は一致しています。")
    else:
        print("consensusValueの値は一致していません。")

# ログファイル1とログファイル2の内容を取得する
log1_content = ""
log2_content = ""

with open("log1.txt", "r") as file:
    log1_content = file.read()

with open("log2.txt", "r") as file:
    log2_content = file.read()

# ログファイルを比較する
compare_consensus_values(log1_content, log2_content)