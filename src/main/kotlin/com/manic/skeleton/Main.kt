package com.manic.skeleton

import com.xenomachina.argparser.ArgParser
import com.xenomachina.argparser.SystemExitException
import com.xenomachina.argparser.mainBody
import java.io.File

fun main(args: Array<String>) = mainBody {
    printSkeletor()
    val firstArgument: Array<out String> = args.takeIf { it.isNotEmpty() }?.sliceArray(0..0) ?: emptyArray()
    ArgParser(firstArgument).parseInto(::SkeletonProgramArgs).run {
        val argsWithoutInitialCommand = args.sliceArray(1 until args.size)
        when (command) {
            "add" -> ArgParser(argsWithoutInitialCommand).parseInto(::AddCommandArgs).run(::handleAdd)
            else -> throw SystemExitException(
                """
                Unkown COMMAND '$command'.
                Use skeleton --help for usable commands.
                """.trimIndent(), 101
            );
        }
    }
}

private fun handleAdd(addCommandArgs: AddCommandArgs) {
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
}

private fun printSkeletor() {
    println(
        """
,____,
|f-"Y\|
\()7L/
 cgD                            __ _
 |\(                          .'  Y '>,
  \ \                        / _   _   \
   \\\                       )(_) (_)(|}
    \\\                      {  4A   } /
     \\\                      \uLuJJ/\l
      \\\                     |3    p)/
       \\\___ __________      /nnm_n//
       c7___-__,__-)\,__)(".  \_>-<_/D
                  //V     \_"-._.__G G_c__.-__<"/ ( \
                         <"-._>__-,G_.___)\   \7\
                        ("-.__.| \"<.__.-" )   \ \
                        |"-.__"\  |"-.__.-".\   \ \
                        ("-.__"". \"-.__.-".|    \_\
                        \"-.__""|!|"-.__.-".)     \ \
                         "-.__""\_|"-.__.-"./      \ l
                          ".__""${'"'}>G>-.__.-">       .--,_
                              ""  G
    """.trimIndent()
    )
}
