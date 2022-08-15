import time
import pymysql
import pymysql.cursors
import datetime

conn = pymysql.connect(
   host='127.0.0.1',
   user='root',
   password='dingyuan841218',
   database='entry_task',
   charset='utf8mb4',
   cursorclass=pymysql.cursors.DictCursor
)

cursor = conn.cursor()

# $2a$04$EMdYwzi3AQH9LpVbI8wg2O9IUfSute3aJVEygGRkyWEN/FXuscz/u is bcrypt encode for "dingyuan"
sql = """INSERT INTO user (usr_name, usr_nickname, usr_password, profile_pic_url, ctime, mtime) 
VALUES
("stress_test_{}","stress_test_nickname_{}","ddf95206809a0520986c09d836b52bdc", "https://c-ssl.duitang.com/uploads/item/202001/27/20200127191847_vfsJ3.jpeg", "{}", "{}");
"""

print(sql)
for i in range(10000000, 20000000):
   try:
      tmp_sql = sql.format(i, i, datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"), datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
      cursor.execute(tmp_sql)
      if not i % 10000:
         print(i)
         print(tmp_sql)
         conn.commit()

   except Exception as e:
      print(e)

# Closing the connection
conn.commit()
conn.close()