# Union Arena 遊戲測試 - Bob vs Kage


### 1. 創建遊戲 (POST /api/v1/games) - 完整 50 張卡組

curl -X 'POST' \
  'http://localhost:8004/api/v1/games' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "player1_id": "94b46616-3b46-41b3-81dc-e95f70bfb7d5",
  "player2_id": "a8e16546-5a86-415a-9baa-ae62b13891b4",
  "game_mode": "casual",
  "player1_deck": [
    {"id": "00000000-0000-0000-0000-000000000001", "card_number": "UA25BT-001", "card_variant_id": "UA25BT-001-C", "name": "Bob的角色卡1", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3000, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡1", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000002", "card_number": "UA25BT-002", "card_variant_id": "UA25BT-002-C", "name": "Bob的角色卡2", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2500, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡2", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000003", "card_number": "UA25BT-003", "card_variant_id": "UA25BT-003-C", "name": "Bob的角色卡3", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3500, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡3", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000004", "card_number": "UA25BT-004", "card_variant_id": "UA25BT-004-R", "name": "Bob的角色卡4", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 4000, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["戦士"], "effect_text": "Bob的角色卡4", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000005", "card_number": "UA25BT-005", "card_variant_id": "UA25BT-005-SR", "name": "Bob的王牌角色", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 5000, "ap_cost": 4, "energy_cost": "{\"red\": 4}", "energy_produce": "{\"red\": 1}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["船長"], "effect_text": "Bob的王牌角色", "trigger_effect": "DRAW_CARD", "keywords": ["突破"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    
    {"id": "00000000-0000-0000-0000-000000000021", "card_number": "UA25BT-021", "card_variant_id": "UA25BT-021-C", "name": "Bob的角色卡6", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2000, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡6", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000022", "card_number": "UA25BT-022", "card_variant_id": "UA25BT-022-C", "name": "Bob的角色卡7", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2800, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡7", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000023", "card_number": "UA25BT-023", "card_variant_id": "UA25BT-023-C", "name": "Bob的角色卡8", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3200, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡8", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000024", "card_number": "UA25BT-024", "card_variant_id": "UA25BT-024-C", "name": "Bob的角色卡9", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3800, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡9", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000025", "card_number": "UA25BT-025", "card_variant_id": "UA25BT-025-C", "name": "Bob的角色卡10", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 4200, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡10", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000026", "card_number": "UA25BT-026", "card_variant_id": "UA25BT-026-C", "name": "Bob的角色卡11", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2200, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡11", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000027", "card_number": "UA25BT-027", "card_variant_id": "UA25BT-027-C", "name": "Bob的角色卡12", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2600, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡12", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000028", "card_number": "UA25BT-028", "card_variant_id": "UA25BT-028-C", "name": "Bob的角色卡13", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3300, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡13", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000029", "card_number": "UA25BT-029", "card_variant_id": "UA25BT-029-C", "name": "Bob的角色卡14", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3700, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡14", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000030", "card_number": "UA25BT-030", "card_variant_id": "UA25BT-030-R", "name": "Bob的角色卡15", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 4100, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["戦士"], "effect_text": "Bob的角色卡15", "trigger_effect": "DRAW_CARD", "keywords": ["突破"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000031", "card_number": "UA25BT-031", "card_variant_id": "UA25BT-031-C", "name": "Bob的角色卡16", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2400, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡16", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000032", "card_number": "UA25BT-032", "card_variant_id": "UA25BT-032-C", "name": "Bob的角色卡17", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2700, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡17", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000033", "card_number": "UA25BT-033", "card_variant_id": "UA25BT-033-C", "name": "Bob的角色卡18", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3400, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡18", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000034", "card_number": "UA25BT-034", "card_variant_id": "UA25BT-034-C", "name": "Bob的角色卡19", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3600, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡19", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000035", "card_number": "UA25BT-035", "card_variant_id": "UA25BT-035-C", "name": "Bob的角色卡20", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3900, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡20", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000006", "card_number": "UA25BT-006", "card_variant_id": "UA25BT-006-C", "name": "Bob的AP卡1", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000007", "card_number": "UA25BT-007", "card_variant_id": "UA25BT-007-C", "name": "Bob的AP卡2", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000008", "card_number": "UA25BT-008", "card_variant_id": "UA25BT-008-C", "name": "Bob的AP卡3", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000009", "card_number": "UA25BT-009", "card_variant_id": "UA25BT-009-C", "name": "Bob的事件卡1", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "攻撃力+1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000036", "card_number": "UA25BT-036", "card_variant_id": "UA25BT-036-C", "name": "Bob的事件卡2", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "防禦力+1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000037", "card_number": "UA25BT-037", "card_variant_id": "UA25BT-037-C", "name": "Bob的事件卡3", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "抽2張卡", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000038", "card_number": "UA25BT-038", "card_variant_id": "UA25BT-038-C", "name": "Bob的事件卡4", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "角色復活", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000039", "card_number": "UA25BT-039", "card_variant_id": "UA25BT-039-C", "name": "Bob的事件卡5", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "敵方角色麻痺", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000040", "card_number": "UA25BT-040", "card_variant_id": "UA25BT-040-C", "name": "Bob的事件卡6", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "直接傷害1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000041", "card_number": "UA25BT-041", "card_variant_id": "UA25BT-041-C", "name": "Bob的事件卡7", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "全體角色+2000BP", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000042", "card_number": "UA25BT-042", "card_variant_id": "UA25BT-042-C", "name": "Bob的事件卡8", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "破壞敵方場域", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000043", "card_number": "UA25BT-043", "card_variant_id": "UA25BT-043-C", "name": "Bob的事件卡9", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "角色移動", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000044", "card_number": "UA25BT-044", "card_variant_id": "UA25BT-044-C", "name": "Bob的事件卡10", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "能量回復", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000045", "card_number": "UA25BT-045", "card_variant_id": "UA25BT-045-C", "name": "Bob的事件卡11", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "連續攻撃", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000046", "card_number": "UA25BT-046", "card_variant_id": "UA25BT-046-C", "name": "Bob的事件卡12", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 0}", "rarity": "R", "rarity_code": "R", "characteristics": ["戦術"], "effect_text": "究極技能", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000047", "card_number": "UA25BT-047", "card_variant_id": "UA25BT-047-C", "name": "Bob的事件卡13", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "反擊", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000048", "card_number": "UA25BT-048", "card_variant_id": "UA25BT-048-C", "name": "Bob的事件卡14", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "戦術變更", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000049", "card_number": "UA25BT-049", "card_variant_id": "UA25BT-049-C", "name": "Bob的事件卡15", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "最後一擊", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000010", "card_number": "UA25BT-010", "card_variant_id": "UA25BT-010-C", "name": "Bob的場域卡1", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["船"], "effect_text": "紅色角色+500BP", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000050", "card_number": "UA25BT-050", "card_variant_id": "UA25BT-050-C", "name": "Bob的場域卡2", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["訓練場"], "effect_text": "每回合抽額外1張卡", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000051", "card_number": "UA25BT-051", "card_variant_id": "UA25BT-051-C", "name": "Bob的場域卡3", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["港口"], "effect_text": "角色召喚費用-1", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000052", "card_number": "UA25BT-052", "card_variant_id": "UA25BT-052-C", "name": "Bob的場域卡4", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 0}", "rarity": "R", "rarity_code": "R", "characteristics": ["要塞"], "effect_text": "全體角色獲得防壁", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000053", "card_number": "UA25BT-053", "card_variant_id": "UA25BT-053-C", "name": "Bob的場域卡5", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["競技場"], "effect_text": "戦闘時BP+1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000054", "card_number": "UA25BT-054", "card_variant_id": "UA25BT-054-C", "name": "Bob的場域卡6", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["海域"], "effect_text": "移動範圍+1", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000055", "card_number": "UA25BT-055", "card_variant_id": "UA25BT-055-C", "name": "Bob的場域卡7", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["武器庫"], "effect_text": "裝備效果+50%", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000056", "card_number": "UA25BT-056", "card_variant_id": "UA25BT-056-C", "name": "Bob的場域卡8", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["島嶼"], "effect_text": "每回合恢復AP+1", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000057", "card_number": "UA25BT-057", "card_variant_id": "UA25BT-057-C", "name": "Bob的場域卡9", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 4, "energy_cost": "{\"red\": 4}", "energy_produce": "{\"red\": 0}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["傳說之地"], "effect_text": "勝利條件追加", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000058", "card_number": "UA25BT-058", "card_variant_id": "UA25BT-058-C", "name": "Bob的場域卡10", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["小船"], "effect_text": "快速移動", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000059", "card_number": "UA25BT-059", "card_variant_id": "UA25BT-059-C", "name": "Bob的補充卡1", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2300, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的補充卡1", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000060", "card_number": "UA25BT-060", "card_variant_id": "UA25BT-060-C", "name": "Bob的補充卡2", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2800, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的補充卡2", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
  ],
  "player2_deck": [
    {"id": "00000000-0000-0000-0000-000000000011", "card_number": "UA25BT-011", "card_variant_id": "UA25BT-011-C", "name": "Kage的角色卡1", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3500, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡1", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000012", "card_number": "UA25BT-012", "card_variant_id": "UA25BT-012-C", "name": "Kage的角色卡2", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3000, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡2", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000013", "card_number": "UA25BT-013", "card_variant_id": "UA25BT-013-C", "name": "Kage的角色卡3", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4000, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡3", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000014", "card_number": "UA25BT-014", "card_variant_id": "UA25BT-014-R", "name": "Kage的角色卡4", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4500, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["忍者"], "effect_text": "Kage的角色卡4", "trigger_effect": "COLOR", "keywords": ["隠密"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000015", "card_number": "UA25BT-015", "card_variant_id": "UA25BT-015-SR", "name": "Kage的王牌角色", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 5500, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 1}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["影"], "effect_text": "Kage的王牌角色", "trigger_effect": "COLOR", "keywords": ["隠密", "突破"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    
    {"id": "00000000-0000-0000-0000-000000000061", "card_number": "UA25BT-061", "card_variant_id": "UA25BT-061-C", "name": "Kage的角色卡6", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 2800, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡6", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000062", "card_number": "UA25BT-062", "card_variant_id": "UA25BT-062-C", "name": "Kage的角色卡7", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3200, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡7", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000063", "card_number": "UA25BT-063", "card_variant_id": "UA25BT-063-C", "name": "Kage的角色卡8", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3800, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡8", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000064", "card_number": "UA25BT-064", "card_variant_id": "UA25BT-064-C", "name": "Kage的角色卡9", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4200, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡9", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000065", "card_number": "UA25BT-065", "card_variant_id": "UA25BT-065-C", "name": "Kage的角色卡10", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4600, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡10", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000066", "card_number": "UA25BT-066", "card_variant_id": "UA25BT-066-C", "name": "Kage的角色卡11", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 2600, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡11", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000067", "card_number": "UA25BT-067", "card_variant_id": "UA25BT-067-C", "name": "Kage的角色卡12", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3100, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡12", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000068", "card_number": "UA25BT-068", "card_variant_id": "UA25BT-068-C", "name": "Kage的角色卡13", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3700, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡13", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000069", "card_number": "UA25BT-069", "card_variant_id": "UA25BT-069-C", "name": "Kage的角色卡14", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4100, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡14", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000070", "card_number": "UA25BT-070", "card_variant_id": "UA25BT-070-R", "name": "Kage的角色卡15", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4800, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["影"], "effect_text": "Kage的角色卡15", "trigger_effect": "COLOR", "keywords": ["隠密"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000071", "card_number": "UA25BT-071", "card_variant_id": "UA25BT-071-C", "name": "Kage的角色卡16", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 2900, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡16", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000072", "card_number": "UA25BT-072", "card_variant_id": "UA25BT-072-C", "name": "Kage的角色卡17", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3300, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡17", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000073", "card_number": "UA25BT-073", "card_variant_id": "UA25BT-073-C", "name": "Kage的角色卡18", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3900, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡18", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000074", "card_number": "UA25BT-074", "card_variant_id": "UA25BT-074-C", "name": "Kage的角色卡19", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4300, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡19", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000075", "card_number": "UA25BT-075", "card_variant_id": "UA25BT-075-C", "name": "Kage的角色卡20", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4700, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["影"], "effect_text": "Kage的角色卡20", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000016", "card_number": "UA25BT-016", "card_variant_id": "UA25BT-016-C", "name": "Kage的AP卡1", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000017", "card_number": "UA25BT-017", "card_variant_id": "UA25BT-017-C", "name": "Kage的AP卡2", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000018", "card_number": "UA25BT-018", "card_variant_id": "UA25BT-018-C", "name": "Kage的AP卡3", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000019", "card_number": "UA25BT-019", "card_variant_id": "UA25BT-019-C", "name": "Kage的事件卡1", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "角色歸位", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000076", "card_number": "UA25BT-076", "card_variant_id": "UA25BT-076-C", "name": "Kage的事件卡2", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "隠身術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000077", "card_number": "UA25BT-077", "card_variant_id": "UA25BT-077-C", "name": "Kage的事件卡3", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "分身術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000078", "card_number": "UA25BT-078", "card_variant_id": "UA25BT-078-C", "name": "Kage的事件卡4", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "反擊準備", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000079", "card_number": "UA25BT-079", "card_variant_id": "UA25BT-079-C", "name": "Kage的事件卡5", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "水遁術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000080", "card_number": "UA25BT-080", "card_variant_id": "UA25BT-080-C", "name": "Kage的事件卡6", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "情報收集", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000081", "card_number": "UA25BT-081", "card_variant_id": "UA25BT-081-C", "name": "Kage的事件卡7", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "影分身大軍", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000082", "card_number": "UA25BT-082", "card_variant_id": "UA25BT-082-C", "name": "Kage的事件卡8", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "戦術撤退", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000083", "card_number": "UA25BT-083", "card_variant_id": "UA25BT-083-C", "name": "Kage的事件卡9", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "瞬身術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000084", "card_number": "UA25BT-084", "card_variant_id": "UA25BT-084-C", "name": "Kage的事件卡10", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "完美配合", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000085", "card_number": "UA25BT-085", "card_variant_id": "UA25BT-085-C", "name": "Kage的事件卡11", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "手裏劍術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000086", "card_number": "UA25BT-086", "card_variant_id": "UA25BT-086-C", "name": "Kage的事件卡12", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 0}", "rarity": "R", "rarity_code": "R", "characteristics": ["忍術"], "effect_text": "奧義・影縛術", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000087", "card_number": "UA25BT-087", "card_variant_id": "UA25BT-087-C", "name": "Kage的事件卡13", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "陷阱設置", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000088", "card_number": "UA25BT-088", "card_variant_id": "UA25BT-088-C", "name": "Kage的事件卡14", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術"], "effect_text": "煙霧彈", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000089", "card_number": "UA25BT-089", "card_variant_id": "UA25BT-089-C", "name": "Kage的事件卡15", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "致命一擊", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},

    {"id": "00000000-0000-0000-0000-000000000020", "card_number": "UA25BT-020", "card_variant_id": "UA25BT-020-C", "name": "Kage的場域卡1", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術道場"], "effect_text": "藍色角色+500BP", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000090", "card_number": "UA25BT-090", "card_variant_id": "UA25BT-090-C", "name": "Kage的場域卡2", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["隠れ家"], "effect_text": "隠密效果強化", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000091", "card_number": "UA25BT-091", "card_variant_id": "UA25BT-091-C", "name": "Kage的場域卡3", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["水池"], "effect_text": "忍術費用-1", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000092", "card_number": "UA25BT-092", "card_variant_id": "UA25BT-092-C", "name": "Kage的場域卡4", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 0}", "rarity": "R", "rarity_code": "R", "characteristics": ["影之森"], "effect_text": "全體角色獲得隠密", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000093", "card_number": "UA25BT-093", "card_variant_id": "UA25BT-093-C", "name": "Kage的場域卡5", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["修練場"], "effect_text": "戦術時BP+1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000094", "card_number": "UA25BT-094", "card_variant_id": "UA25BT-094-C", "name": "Kage的場域卡6", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["屋頂"], "effect_text": "移動距離+2", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000095", "card_number": "UA25BT-095", "card_variant_id": "UA25BT-095-C", "name": "Kage的場域卡7", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["武器庫"], "effect_text": "忍具效果+50%", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000096", "card_number": "UA25BT-096", "card_variant_id": "UA25BT-096-C", "name": "Kage的場域卡8", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["月夜"], "effect_text": "夜間戦闘強化", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000097", "card_number": "UA25BT-097", "card_variant_id": "UA25BT-097-C", "name": "Kage的場域卡9", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 0}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["影の國"], "effect_text": "究極忍術解鎖", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000098", "card_number": "UA25BT-098", "card_variant_id": "UA25BT-098-C", "name": "Kage的場域卡10", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["樹枝"], "effect_text": "快速逃脫", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000099", "card_number": "UA25BT-099", "card_variant_id": "UA25BT-099-C", "name": "Kage的補充卡1", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 2700, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的補充卡1", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000100", "card_number": "UA25BT-100", "card_variant_id": "UA25BT-100-C", "name": "Kage的補充卡2", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3100, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的補充卡2", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
  ]
}'

## 2.登入取得jwt token 用戶信息

**Player 1 (Bob):**
- User ID: `94b46616-3b46-41b3-81dc-e95f70bfb7d5`
- Username: `bob`
- Password: `bobbob`

```bash
curl -X 'POST' \
  'http://localhost:8002/api/v1/auth/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "identifier": "bob",
  "password": "bobbob"
}'
```

**Player 2 (Kage):**
- User ID: `a8e16546-5a86-415a-9baa-ae62b13891b4`
- Username: `kage`
- Password: `kagekage`

```bash
curl -X 'POST' \
  'http://localhost:8002/api/v1/auth/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "identifier": "kage",
  "password": "kagekage"
}'
```


### 3. 加入遊戲 (Join Game)
記錄創建遊戲後返回的 `game_id`，然後：

#### Bob 和 Kage 分別加入遊戲

**Player 1 (Bob) 加入遊戲:**
```bash
curl -X 'POST' \
  'http://localhost:8004/api/v1/games/{gameId}/join' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer {bob_token}' \
  -d ''
```

**Player 2 (Kage) 加入遊戲:**
```bash
curl -X 'POST' \
  'http://localhost:8004/api/v1/games/{gameId}/join' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer {kage_token}' \
  -d ''
```

### 4. 檢查遊戲狀態 (Check Game Status)
```bash
curl -X 'GET' \
  'http://localhost:8004/api/v1/game-info/{gameId}' \
  -H 'accept: application/json'
```

確認返回的 `"status": "IN_PROGRESS"` 才能繼續下一步。


### 5. 調度階段 (Mulligan Phase)
Bob 和 Kage 分別決定是否重抽手牌：

**Bob 的調度決定:**
```bash
curl -X 'POST' \
  'http://localhost:8004/api/v1/games/{gameId}/mulligan' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer {bob_token}' \
  -d '{
    "mulligan": false
  }'
```

**Kage 的調度決定:**
```bash
curl -X 'POST' \
  'http://localhost:8004/api/v1/games/{gameId}/mulligan' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer {kage_token}' \
  -d '{
    "mulligan": false
  }'
```

> `mulligan` 設為 `true` 表示重抽，`false` 表示保留當前手牌

### 6. 查詢完整遊戲狀態  
調用 `GET /api/v1/games/{gameId}` 查看完整遊戲狀態：

```bash
curl -X 'GET' \
  'http://localhost:8004/api/v1/games/{gameId}' \
  -H 'Authorization: Bearer {bob_token}'
```

#### 執行遊戲動作

**重要：使用新的 turn-info API 確定當前回合玩家**

首先查詢當前回合信息：
```bash
curl -X GET "http://localhost:8004/api/v1/games/{gameId}/turn-info"
```

根據回應中的 `is_player1_turn` 和 `is_player2_turn` 決定使用哪個玩家的 JWT token。

**測試所有遊戲動作：**

#### 1. 抽牌動作 (DRAW_CARD)
```bash
# Bob的回合時使用 (Player1)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "DRAW_CARD",
    "action_data": []
  }'

# Kage的回合時使用 (Player2) 
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "DRAW_CARD",
    "action_data": []
  }'
```

#### 2. 額外抽牌動作 (EXTRA_DRAW) - 支付1AP額外抽1張
```bash
# Bob的回合時使用 (Player1) - 起始階段可用
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "EXTRA_DRAW",
    "action_data": []
  }'

# Kage的回合時使用 (Player2) - 起始階段可用
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "EXTRA_DRAW",
    "action_data": []
  }'
```

#### 3. 出牌動作 (PLAY_CARD) - 主要階段可用
```bash
# 出角色卡到能源線 (Bob的回合)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000001",
      "position": {
        "zone": "energy_line",
        "index": 0
      }
    }
  }'

# 出角色卡到前線 (Bob的回合)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000001",
      "position": {
        "zone": "front_line",
        "index": 0
      }
    }
  }'

# 出場域卡 (Bob的回合)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000010"
    }
  }'

# 出事件卡 (Bob的回合)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000009"
    }
  }'

# 出AP卡 (提供能源)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000006"
    }
  }'

# Kage出牌範例
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "PLAY_CARD",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000011",
      "position": {
        "zone": "energy_line",
        "index": 0
      }
    }
  }'
```

#### 4. 攻擊動作 (ATTACK) - 攻擊階段可用
```bash
# 攻擊對手玩家 (Bob攻擊Kage)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "ATTACK",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000001",
      "target_type": "player"
    }
  }'

# 攻擊對手角色卡 (Bob攻擊Kage的角色)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "ATTACK",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000001",
      "target_type": "character",
      "target_id": "00000000-0000-0000-0000-000000000011"
    }
  }'

# Kage攻擊Bob
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "ATTACK",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000011",
      "target_type": "player"
    }
  }'
```

#### 5. 角色移動動作 (MOVE_CHARACTER) - 移動階段可用
```bash
# 從能源線移動到前線 (Bob)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "MOVE_CHARACTER",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000001",
      "from_zone": "energy_line",
      "to_zone": "front_line",
      "to_index": 0
    }
  }'

# 前線內位置調整 (Kage)
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "MOVE_CHARACTER",
    "action_data": {
      "card_id": "00000000-0000-0000-0000-000000000011",
      "from_zone": "front_line",
      "to_zone": "front_line",
      "to_index": 1
    }
  }'
```

#### 6. 結束階段動作 (END_PHASE)
```bash
# Bob結束當前階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "END_PHASE",
    "action_data": []
  }'

# Kage結束當前階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "END_PHASE",
    "action_data": []
  }'
```

#### 7. 結束回合動作 (END_TURN)
```bash
# Bob結束回合
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "END_TURN",
    "action_data": []
  }'

# Kage結束回合
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "END_TURN",
    "action_data": []
  }'
```

#### 8. 投降動作 (SURRENDER)
```bash
# Bob投降
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "SURRENDER",
    "action_data": []
  }'

# Kage投降
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "SURRENDER",
    "action_data": []
  }'
```

## 🔍 測試錯誤情況

#### 測試非當前回合玩家執行動作 (應返回 403 Forbidden)
```bash
# 當前是Bob回合時，Kage嘗試執行動作
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {kage_token}" \
  -d '{
    "action_type": "DRAW_CARD",
    "action_data": []
  }' \
  -w "\nHTTP Status: %{http_code}\n"
```

#### 測試無效動作類型 (應返回 500 Internal Server Error)
```bash
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "INVALID_ACTION",
    "action_data": []
  }' \
  -w "\nHTTP Status: %{http_code}\n"
```

#### 測試不存在的遊戲ID (應返回 404 Not Found)
```bash
curl -X POST "http://localhost:8004/api/v1/games/00000000-0000-0000-0000-000000000000/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{
    "action_type": "DRAW_CARD",
    "action_data": []
  }' \
  -w "\nHTTP Status: %{http_code}\n"
```

## 📝 完整遊戲流程測試範例

```bash
# 1. 確認當前回合
curl -X GET "http://localhost:8004/api/v1/games/{gameId}/turn-info"

# 2. 如果是Bob回合且在起始階段 - 抽牌或額外抽牌
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "DRAW_CARD", "action_data": []}'

# 3. 結束起始階段，進入移動階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "END_PHASE", "action_data": []}'

# 4. 移動階段 - 移動角色位置
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "MOVE_CHARACTER", "action_data": {"card_id": "00000000-0000-0000-0000-000000000001", "from_zone": "energy_line", "to_zone": "front_line", "to_index": 0}}'

# 5. 結束移動階段，進入主要階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "END_PHASE", "action_data": []}'

# 6. 主要階段 - 出牌
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "PLAY_CARD", "action_data": {"card_id": "00000000-0000-0000-0000-000000000002", "position": {"zone": "energy_line", "index": 1}}}'

# 7. 結束主要階段，進入攻擊階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "END_PHASE", "action_data": []}'

# 8. 攻擊階段 - 攻擊
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "ATTACK", "action_data": {"card_id": "00000000-0000-0000-0000-000000000001", "target_type": "player"}}'

# 9. 結束攻擊階段，進入結束階段
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "END_PHASE", "action_data": []}'

# 10. 結束回合
curl -X POST "http://localhost:8004/api/v1/games/{gameId}/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {bob_token}" \
  -d '{"action_type": "END_TURN", "action_data": []}'
```

## 🎯 JWT Token 變量替換

記得將以下變量替換為實際值：
- `{gameId}`: 實際的遊戲ID
- `{bob_token}`: Bob登錄後獲得的access_token  
- `{kage_token}`: Kage登錄後獲得的access_token

## 🎮 命令行快速測試

### 完整測試流程 - 使用實際數據

假設你已經有以下實際數據，執行完整測試流程：

```bash
# 設置變量
export GAME_ID="your-game-id-here"
export BOB_TOKEN="your-bob-token-here"
export KAGE_TOKEN="your-kage-token-here"

# 1. 檢查當前回合
curl -X GET "http://localhost:8004/api/v1/games/$GAME_ID/turn-info"

# 2. Bob 執行動作（如果是 Bob 的回合）
curl -X POST "http://localhost:8004/api/v1/games/$GAME_ID/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{"action_type": "DRAW_CARD", "action_data": []}'

# 3. 查看遊戲狀態變化
curl -X GET "http://localhost:8004/api/v1/games/$GAME_ID" \
  -H "Authorization: Bearer $BOB_TOKEN"
```

### 錯誤處理測試

```bash
# 測試非當前回合玩家執行動作（應該返回 403）
curl -X POST "http://localhost:8004/api/v1/games/$GAME_ID/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $KAGE_TOKEN" \
  -d '{"action_type": "DRAW_CARD", "action_data": []}' \
  -w "\nHTTP Status: %{http_code}\n"

# 測試無效動作類型（應該返回 500）
curl -X POST "http://localhost:8004/api/v1/games/$GAME_ID/actions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{"action_type": "INVALID_ACTION", "action_data": []}' \
  -w "\nHTTP Status: %{http_code}\n"
```

## 📋 測試檢查清單

### 基本流程測試
- [ ] 創建遊戲（50張卡組）
- [ ] 兩玩家登錄獲取 JWT token
- [ ] 兩玩家加入遊戲
- [ ] 檢查遊戲狀態變為 IN_PROGRESS
- [ ] 調度階段（兩玩家決定是否重抽）
- [ ] 查詢 turn-info 確認當前回合玩家
- [ ] 執行各種遊戲動作

### 錯誤測試
- [ ] 非當前回合玩家執行動作 → 403 Forbidden
- [ ] 無效動作類型 → 500 Internal Server Error
- [ ] 不存在的遊戲ID → 404 Not Found
- [ ] 無效的 JWT token → 401 Unauthorized

### 動作類型測試
- [ ] DRAW_CARD - 抽牌
- [ ] EXTRA_DRAW - 額外抽牌
- [ ] PLAY_CARD - 出牌（角色卡、場域卡、事件卡、AP卡）
- [ ] MOVE_CHARACTER - 角色移動
- [ ] ATTACK - 攻擊
- [ ] END_PHASE - 結束階段
- [ ] END_TURN - 結束回合
- [ ] SURRENDER - 投降

## 🔧 常見問題排除

1. **Token 過期**: 如果收到 401 錯誤，重新登錄獲取新的 access_token
2. **遊戲不存在**: 確認 gameId 正確，檢查遊戲是否已經結束
3. **卡組驗證失敗**: 確保使用 50 張卡片的完整卡組
4. **動作無效**: 檢查當前遊戲階段是否允許該動作類型

---

**📝 注意**: 記得將所有 `{gameId}`, `{bob_token}`, `{kage_token}` 替換為實際值！
