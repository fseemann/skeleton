<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>${groupId}</groupId>
        <artifactId>${artifactId}</artifactId>
        <version>${version}</version>
    </parent>

    <artifactId>${artifactId}-infrastructure</artifactId>
    <version>${version}</version>

    <dependencies>
        <dependency>
            <groupId>${groupId}</groupId>
            <artifactId>${artifactId}-application</artifactId>
            <version>${version}</version>
        </dependency>
    </dependencies>
</project>