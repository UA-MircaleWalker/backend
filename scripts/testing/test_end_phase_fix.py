#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
測試END_PHASE修復
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
import time

# 設置
BASE_URL = "http://localhost:8004/api/v1"
USER_SERVICE_URL = "http://localhost:8002/api/v1"
BOB_ID = "94b46616-3b46-41b3-81dc-e95f70bfb7d5"

def login_user():
    """登入用戶並返回token"""
    url = f"{USER_SERVICE_URL}/auth/login"
    data = {"identifier": "bob", "password": "bobbob"}
    response = requests.post(url, json=data)
    
    if response.status_code == 200:
        result = response.json()
        if 'data' in result and 'access_token' in result['data']:
            return result['data']['access_token']
        elif 'access_token' in result:
            return result['access_token']
    return None

def create_game():
    """創建遊戲"""
    with open("C:\\Users\\weilo\\Desktop\\ua\\test_data\\FULL_50_CARDS_DECK.json", "r", encoding="utf-8") as f:
        game_data = json.load(f)
    
    response = requests.post(f"{BASE_URL}/games", json=game_data)
    if response.status_code == 201:
        result = response.json()
        if 'data' in result and 'game' in result['data']:
            return result['data']['game']['id']
        elif 'game' in result:
            return result['game']['id']
    return None

def join_game(game_id, token):
    """加入遊戲"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    response = requests.post(f"{BASE_URL}/games/{game_id}/join", headers=headers)
    return response.status_code == 200

def perform_action(game_id, action_type, token):
    """執行遊戲動作"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    data = {
        "game_id": game_id,
        "player_id": BOB_ID,
        "action_type": action_type,
        "action_data": []
    }
    
    print(f"🎯 執行動作: {action_type}")
    response = requests.post(f"{BASE_URL}/games/{game_id}/actions", headers=headers, json=data)
    print(f"📊 狀態碼: {response.status_code}")
    
    if response.status_code == 200:
        result = response.json()
        print(f"✅ 成功: {result.get('success', False)}")
        if not result.get('success'):
            print(f"❌ 錯誤: {result.get('error', 'Unknown')}")
        return result
    else:
        try:
            error_data = response.json()
            print(f"❌ 失敗: {error_data.get('error', 'Unknown error')}")
        except:
            print(f"❌ 失敗: HTTP {response.status_code}")
        return {"success": False, "error": f"HTTP {response.status_code}"}

def main():
    print("🚀 測試 END_PHASE 修復")
    
    # 1. 登入
    print("\n1. 登入...")
    token = login_user()
    if not token:
        print("❌ 登入失敗")
        return
    print(f"✅ 登入成功")
    
    # 2. 創建遊戲
    print("\n2. 創建遊戲...")
    game_id = create_game()
    if not game_id:
        print("❌ 創建遊戲失敗")
        return
    print(f"✅ 遊戲創建成功: {game_id}")
    
    # 3. 加入遊戲
    print("\n3. 加入遊戲...")
    if not join_game(game_id, token):
        print("❌ 加入遊戲失敗")
        return
    print("✅ 加入遊戲成功")
    
    # 4. 執行調度
    print("\n4. 執行調度...")
    mulligan_data = {
        "game_id": game_id,
        "player_id": BOB_ID,
        "mulligan": False
    }
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    response = requests.post(f"{BASE_URL}/games/{game_id}/mulligan", headers=headers, json=mulligan_data)
    print(f"📊 調度狀態碼: {response.status_code}")
    
    time.sleep(1)
    
    # 5. 測試END_PHASE
    print("\n5. 測試 END_PHASE...")
    result = perform_action(game_id, "END_PHASE", token)
    
    if result.get('success'):
        print("🎉 END_PHASE 修復成功！")
    else:
        print("❌ END_PHASE 仍然失敗")
        
    print(f"\n📋 測試完成")

if __name__ == "__main__":
    main()