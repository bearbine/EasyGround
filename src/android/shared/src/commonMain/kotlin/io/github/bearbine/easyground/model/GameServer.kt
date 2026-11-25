package io.github.bearbine.easyground.model

// игровой сервер, размещенный на хосте
//описание сервера берем в дату для удобств по типу: copy()
//так же сравнение по данным где они считаются одинаковыми
data class GameServer(
    val id:String, // уник ид
    val hostId:String,//ид хоста на котором он лежит
    val name:String,//нейм(от юзера)
    val gameTemplateId:String,//ид шаблона
    val memoryLimitMiB:Int,//выделенный лимит памяти(потом поменяем, пока что для тестов)
    val observedState:ServerState = ServerState.UNKNOWN, // последнее состояние, по умол стоит unknown
)
