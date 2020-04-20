package xyz.hoper.utils.vertx

interface BaseAsyncService {
    val address: String?
        get() {
            val className = this.javaClass.name
            return className.substring(0, className.lastIndexOf("Impl")).replace(".impl", "")
        }

    @get:Throws(ClassNotFoundException::class)
    val asyncInterfaceClass: Class<*>?
        get() {
            val className = this.javaClass.name
            return Class.forName(className.substring(0, className.lastIndexOf("Impl")).replace(".impl", ""))
        }
}