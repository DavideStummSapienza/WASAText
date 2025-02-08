<template>
  <div class="chat-card">
    <img :src="chat.avatar" alt="Avatar" class="avatar" />
    <div class="chat-info">
      <strong>{{ chat.username }}</strong>
      <p>{{ chat.lastMessage }}</p>
    </div>
    <span class="chat-time">{{ chat.timestamp }}</span>
  </div>
</template>

<script>
export default {
  props: {
    chat: Object
  }
};
</script>

<style scoped>
.chat-card {
  display: flex;
  align-items: center;
  background: white;
  border-radius: 10px;
  padding: 10px;
  margin-bottom: 10px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.chat-info {
  flex-grow: 1;
}

.chat-time {
  font-size: 12px;
  color: gray;
}
</style>