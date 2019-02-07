package com.manic.skeleton

import freemarker.template.Configuration
import freemarker.template.TemplateExceptionHandler

object FreemarkerConfig : Configuration(Configuration.VERSION_2_3_28) {
    init {
        setClassForTemplateLoading(this.javaClass, "/templates")
        defaultEncoding = "UTF-8"
        FreemarkerConfig.templateExceptionHandler = TemplateExceptionHandler.RETHROW_HANDLER
        logTemplateExceptions = false
    }
}