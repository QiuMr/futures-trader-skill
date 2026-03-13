import dataset

_db = dataset.connect(
    "mysql+pymysql://testt1:yPT4eBmg3PW55tKp@mysql5.sqlpub.com:3310/testt1?charset=utf8mb4"
)

with _db as _tx:
    _table = _tx["api_keys"]
    
    invalid_keys = [
        "0e8f8d55ab6f4c91f47b2807bf194388",
        "d160c47603afc534b873c526e82cfbf7"
    ]
    
    for key in invalid_keys:
        _table.delete(key_val=key)

_db.close()

print("已删除验证失败的账号")
