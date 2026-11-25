package io.github.bearbine.easyground

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform