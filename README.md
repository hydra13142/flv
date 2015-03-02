flv
===

simple flv parse


New函数创建一个flv对象，其ReadFrom方法读取一个flv文件，WriteTo方法写入文件flv数据；Video和Audio方法说明该flv文件是否有视频/音频数据；Clip方法可以切割文件，Append方法可以连接文件。

一个flv对象包含多个标签（Tag），标签有Script、Audio、Video三个方法以确定标签类型，Keyframe方法表示是否关键帧；Size、SetSize、Time、SetTime方法可以获取/设置标签的大小、时间戳；AudioType和VideoType方法给出音视频的详细编码和帧类型。

flv对象的Meta字段为MetaData类型，记录了该flv文件的元数据。