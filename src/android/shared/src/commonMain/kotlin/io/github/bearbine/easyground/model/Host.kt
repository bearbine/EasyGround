package io.github.bearbine.easyground.model

//==================
//хост на котором стоят наши сервера
//описание хоста берем в дату для удобств по типу: copy()
//так же сравнение по данным где они считаются одинаковыми
data class Host(
    val id:String,//уник ид хоста
    val name:String,//нейм(от юзера)
)