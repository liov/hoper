name := "akka-grpc-kubernetes"
scalaVersion := "2.13.0"

lazy val akkaVersion = "2.5.23"
lazy val discoveryVersion = "1.0.2"
lazy val akkaHttpVersion = "10.1.9"
lazy val alpnVersion = "2.0.9"

lazy val compileSettings =
  inConfig(Compile)(Seq(
    excludeFilter in PB.generate := "descriptor.proto",
    PB.protoSources += baseDirectory.value / "protobuf",
    PB.protoSources += file("protobuf")
  ))

lazy val commonSettings = Seq(
  version := "0.1.0",
)

lazy val root = (project in file("."))
  .aggregate(httpToGrpc, grpcService)

// Http front end that calls out to a gRPC back end
lazy val httpToGrpc = (project in file("http-to-grpc"))
  .enablePlugins(AkkaGrpcPlugin, DockerPlugin, JavaAppPackaging, JavaAgent)
  .settings(
    commonSettings,
    compileSettings,
    libraryDependencies ++= Seq(
      "com.typesafe.akka" %% "akka-actor" % akkaVersion,
      "com.typesafe.akka" %% "akka-discovery" % akkaVersion,
      "com.typesafe.akka" %% "akka-protobuf" % akkaVersion,
      "com.typesafe.akka" %% "akka-stream" % akkaVersion,

      "com.typesafe.akka" %% "akka-parsing" % akkaHttpVersion,
      "com.typesafe.akka" %% "akka-http-core" % akkaHttpVersion,
      "com.typesafe.akka" %% "akka-http" % akkaHttpVersion,
      "com.typesafe.akka" %% "akka-http-spray-json" % akkaHttpVersion,
      "com.typesafe.akka" %% "akka-http2-support" % akkaHttpVersion,

      "com.lightbend.akka.discovery" %% "akka-discovery-kubernetes-api" % discoveryVersion,
    ),
    javaAgents += "org.mortbay.jetty.alpn" % "jetty-alpn-agent" % alpnVersion % "runtime",
    dockerExposedPorts := Seq(8080),
  )

lazy val grpcService = (project in file("grpc-service"))
  .enablePlugins(AkkaGrpcPlugin, DockerPlugin, JavaAppPackaging, JavaAgent)
  .settings(
    commonSettings,
    compileSettings,
    javaAgents += "org.mortbay.jetty.alpn" % "jetty-alpn-agent" % alpnVersion % "runtime",
    dockerExposedPorts := Seq(8080),
  )
