-- Chat Assistant Backend Database Schema
-- PostgreSQL 表创建脚本
-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- 1. 用户表 (users)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    avatar VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- 创建用户名索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
-- 2. 对话表 (conversations)
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(500) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- 创建对话表索引
CREATE INDEX idx_conversations_user_id ON conversations(user_id);
CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);
CREATE INDEX idx_conversations_provider ON conversations(provider);
-- 3. 消息表 (messages)
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- 创建消息表索引
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_deleted_at ON messages(deleted_at);
CREATE INDEX idx_messages_role ON messages(role);
-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';
-- 为所有表创建更新时间触发器
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_conversations_updated_at BEFORE
UPDATE ON conversations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_messages_updated_at BEFORE
UPDATE ON messages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- 添加表注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON TABLE conversations IS '对话表';
COMMENT ON TABLE messages IS '消息表';
-- 添加字段注释
COMMENT ON COLUMN users.username IS '用户名，唯一';
COMMENT ON COLUMN users.avatar IS '用户头像URL';
COMMENT ON COLUMN conversations.user_id IS '用户ID';
COMMENT ON COLUMN conversations.title IS '对话标题';
COMMENT ON COLUMN conversations.provider IS 'AI提供商 (openai, gemini, local等)';
COMMENT ON COLUMN conversations.model IS 'AI模型 (gpt-4, gemini-pro, llama-3等)';
COMMENT ON COLUMN messages.conversation_id IS '对话ID';
COMMENT ON COLUMN messages.role IS '消息角色 (user, assistant, system)';
COMMENT ON COLUMN messages.content IS '消息内容';