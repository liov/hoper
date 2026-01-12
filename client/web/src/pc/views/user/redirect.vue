<template>
  <div />
</template>

<script>
export default {
  mounted() {
    const token = localStorage.getItem('token')
    if (token) {
      this.$store.commit('SET_TOKEN', token)
      this.$axios.$get(`/api/user/get`).then(res => {
        if (res.code === 200) this.$store.commit('SET_USER', res.data)
         this.$router.replace(this.$route.query.callbackUrl)
      })
    }
  }
}
</script>

<style scoped></style>
