-- 删除所有表的脚本
-- 注意：此脚本会删除所有数据，请谨慎使用！
-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_conversations_updated_at ON conversations;
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
-- 删除函数
DROP FUNCTION IF EXISTS update_updated_at_column();
-- 删除表（按依赖关系顺序）
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- 删除扩展（可选）
-- DROP EXTENSION IF EXISTS "uuid-ossp";