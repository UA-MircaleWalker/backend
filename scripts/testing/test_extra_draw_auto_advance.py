#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
測試 EXTRA_DRAW 自動推進階段功能
"""

import sys
import os

# 設置 UTF-8 編碼輸出
if sys.platform.startswith('win'):
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer, 'strict')
    sys.stderr = codecs.getwriter('utf-8')(sys.stderr.buffer, 'strict')
    os.environ['PYTHONIOENCODING'] = 'utf-8'

import requests
import json

# 設置
BASE_URL = "http://localhost:8004/api/v1"
USER_SERVICE_URL = "http://localhost:8002/api/v1"
BOB_ID = "94b46616-3b46-41b3-81dc-e95f70bfb7d5"

def get_game_status(game_id, token):
    """獲取遊戲狀態"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    
    # 獲取turn-info
    response = requests.get(f"{BASE_URL}/games/{game_id}/turn-info")
    turn_info = response.json()
    
    return turn_info

def perform_extra_draw(game_id, token):
    """執行EXTRA_DRAW動作"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    data = {
        "game_id": game_id,
        "player_id": BOB_ID,
        "action_type": "EXTRA_DRAW",
        "action_data": []
    }
    
    response = requests.post(f"{BASE_URL}/games/{game_id}/actions", headers=headers, json=data)
    return response.json()

def main():
    print("🧪 測試 EXTRA_DRAW 自動階段推進功能")
    
    # 1. 登入
    login_data = {"identifier": "bob", "password": "bobbob"}
    login_response = requests.post(f"{USER_SERVICE_URL}/auth/login", json=login_data)
    token = login_response.json()['data']['access_token']
    
    # 2. 創建遊戲
    with open("C:\\Users\\weilo\\Desktop\\ua\\test_data\\FULL_50_CARDS_DECK.json", "r", encoding="utf-8") as f:
        game_data = json.load(f)
    
    create_response = requests.post(f"{BASE_URL}/games", json=game_data)
    game_id = create_response.json()['data']['game']['id']
    print(f"📄 創建遊戲: {game_id}")
    
    # 3. 加入遊戲
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    join_response = requests.post(f"{BASE_URL}/games/{game_id}/join", headers=headers)
    print(f"🎯 加入遊戲: {'成功' if join_response.status_code == 200 else '失敗'}")
    
    # 4. 調度
    mulligan_data = {"game_id": game_id, "player_id": BOB_ID, "mulligan": False}
    mulligan_response = requests.post(f"{BASE_URL}/games/{game_id}/mulligan", headers=headers, json=mulligan_data)
    print(f"🔄 調度: {'成功' if mulligan_response.status_code == 200 else '失敗'}")
    
    # 5. 檢查初始狀態
    print(f"\n📊 執行 EXTRA_DRAW 前的狀態:")
    turn_info_before = get_game_status(game_id, token)
    
    # 正確解析響應結構
    actual_turn_info_before = turn_info_before['data'] if 'data' in turn_info_before else turn_info_before
    
    print(f"   回合: {actual_turn_info_before.get('turn', 'N/A')}")
    print(f"   階段: {actual_turn_info_before.get('phase', 'N/A')}")
    print(f"   活躍玩家: {'Bob' if actual_turn_info_before.get('is_player1_turn') else 'Kage'}")
    
    # 6. 執行EXTRA_DRAW
    print(f"\n🎯 執行 EXTRA_DRAW...")
    extra_draw_result = perform_extra_draw(game_id, token)
    print(f"   結果: {'成功' if extra_draw_result.get('success') else '失敗'}")
    if extra_draw_result.get('error'):
        print(f"   錯誤: {extra_draw_result['error']}")
    
    # 7. 檢查執行後狀態
    print(f"\n📊 執行 EXTRA_DRAW 後的狀態:")
    turn_info_after = get_game_status(game_id, token)
    
    # 正確解析響應結構
    actual_turn_info_after = turn_info_after['data'] if 'data' in turn_info_after else turn_info_after
    
    print(f"   回合: {actual_turn_info_after.get('turn', 'N/A')}")
    print(f"   階段: {actual_turn_info_after.get('phase', 'N/A')}")
    print(f"   活躍玩家: {'Bob' if actual_turn_info_after.get('is_player1_turn') else 'Kage'}")
    
    # 8. 分析結果
    print(f"\n📋 分析結果:")
    turn_before = actual_turn_info_before.get('turn', 0)
    phase_before = actual_turn_info_before.get('phase', 'UNKNOWN')
    turn_after = actual_turn_info_after.get('turn', 0)
    phase_after = actual_turn_info_after.get('phase', 'UNKNOWN')
    
    print(f"   變化: turn {turn_before} → {turn_after}, phase {phase_before} → {phase_after}")
    
    if turn_before == 1 and phase_before == 0:  # START phase
        if turn_after == 1 and phase_after == 1:  # MOVE phase
            print(f"   ✅ 成功: EXTRA_DRAW 在 START 階段自動推進到 MOVE 階段")
        else:
            print(f"   ❌ 失敗: 應該從 START(0) 推進到 MOVE(1)")
    else:
        print(f"   ⚠️  測試條件不符: 期望在 START 階段執行")

if __name__ == "__main__":
    main()