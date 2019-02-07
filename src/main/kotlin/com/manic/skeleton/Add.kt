package com.manic.skeleton

import com.xenomachina.argparser.SystemExitException
import java.io.File
import java.io.FileWriter

fun add(addCommandArgs: AddCommandArgs) {
    val domainName = addCommandArgs.domainName
    val parent = File(domainName)

    if (parent.exists()) {
        throw SystemExitException("Directory ${parent.absoluteFile} is already existing.", 101)
    }

    val application = File("$domainName/$domainName-application")
    val domain = File("$domainName/$domainName-domain")
    val infrastructure = File("$domainName/$domainName-infrastructure")

    arrayOf(
        parent,
        application,
        domain,
        infrastructure
    ).forEach { it.mkdir() }

    var groupId: String?
    do {
        print("Type group id: ")
        groupId = readLine() ?: ""
        print("Group id '$groupId' correct?[y/n]: ")
    } while(readLine() != "y")

    val parentPom = FreemarkerConfig.getTemplate("add/parent-pom.ftlh")
    parentPom.process(mapOf(
        "groupId" to groupId,
        "domainName" to domainName
    ), FileWriter(File(parent, "pom.xml")))
}