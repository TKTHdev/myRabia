#!/bin/bash

# コンパイルされたGoプログラムの名前
PROGRAM="./GClient"

# 固定パラメータ
DURATION=20

# クライアント数の配列
CLIENT_NUMBERS=(1 2 4 8 16 32 64)

# イテレーションの数
ITERATION_NUMBER=5

# コマンドの配列
COMMANDS=("A" "B" "C")

# 各コマンドに対して処理を実行
for cmd in "${COMMANDS[@]}"; do
    # 出力ファイル
    OUTPUT_FILE="output_${cmd}.txt"

    # ファイルの初期化
    > $OUTPUT_FILE

    echo "Testing Command $cmd" >> $OUTPUT_FILE

    # 各クライアント数で処理
    for clients in "${CLIENT_NUMBERS[@]}"; do
        total_commands=0
        total_read_time=0
        total_write_time=0
        total_total_time=0
        
        echo "Testing with $clients clients" >> $OUTPUT_FILE
        
        # ITERATION_NUMBER回実行してデータを収集
        for ((i=1; i<=ITERATION_NUMBER; i++)); do
            echo "  Run $i of $ITERATION_NUMBER" >> $OUTPUT_FILE
            
            # プログラムを実行し、出力を取得
            temp_output=$(echo -e "$cmd\n$clients\n$DURATION" | $PROGRAM)
            
            # 必要なデータを抽出
            commands=$(echo "$temp_output" | grep "Total number of commands executed:" | awk '{print $NF}')
            read_time=$(echo "$temp_output" | grep "Total Read time:" | awk '{print $NF}')
            write_time=$(echo "$temp_output" | grep "Total write time:" | awk '{print $NF}') 
            total_time=$(echo "$temp_output" | grep "Average total time:" | awk '{print $NF}')
            
            total_commands=$((total_commands + commands))
            total_read_time=$(echo "scale=2; $total_read_time + $read_time" | bc)
            total_write_time=$(echo "scale=2; $total_write_time + $write_time" | bc) 
            total_total_time=$(echo "scale=2; $total_total_time + $total_time" | bc)
            
            echo "    Commands executed: $commands" >> $OUTPUT_FILE
            echo "    Total read latency: $read_time" >> $OUTPUT_FILE
            echo "    Total write latency: $write_time" >> $OUTPUT_FILE
            echo "    Total latency: $total_time" >> $OUTPUT_FILE

            sleep 5  # 各実行間に5秒の休止
        done
        
        # コマンド/秒の平均を計算 
        average_commands=$(echo "scale=2; $total_commands / $ITERATION_NUMBER" | bc)
        commands_per_sec=$(echo "scale=2; $average_commands / $DURATION" | bc)
        
        # 平均レイテンシを計算
        average_read_latency=$(echo "scale=2; $total_read_time / $ITERATION_NUMBER" | bc)
        average_write_latency=$(echo "scale=2; $total_write_time / $ITERATION_NUMBER" | bc)
        average_total_latency=$(echo "scale=2; $total_total_time / $ITERATION_NUMBER" | bc)

        # 結果をファイルに出力
        echo "Clients: $clients, Average commands per second: $commands_per_sec" >> $OUTPUT_FILE
        echo "Clients: $clients, Average read latency: $average_read_latency" >> $OUTPUT_FILE
        echo "Clients: $clients, Average write latency: $average_write_latency" >> $OUTPUT_FILE 
        echo "Clients: $clients, Average total latency: $average_total_latency" >> $OUTPUT_FILE
        echo "----------------------------------------" >> $OUTPUT_FILE 
    done

    echo "Results for Command $cmd have been saved to $OUTPUT_FILE"
    echo "========================================"
done