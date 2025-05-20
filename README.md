# Miniproject go-hex

## About Project
โปรเจ็คนี้สร้างขึ้นเพื่อการทดสอบความรู้ความเข้าใจในการเขียนภาษา golang โดยมีโครงสร้างตาม hexagonal architecture (port and adapter) นำมาลองทำเป็นการสั่งสินค้าผ่าน api

## Work List
* Authentication with jwt 
* Authorization with jwt
* Database mongodb, postgresql  
* Cache by Redis
* Unit Test

## Project Detail
* entities ของโปรเจ็คจะประกอบด้วย 3 entitie ได้เเก่ user, product, order โดย entities เหล่านี้จะเป็นข้อมูลตัวกลางในการพูดคุยระหว่าง service เเละจะถูกนำไปเเปลงอีกทีเมื่อมีการเขียนลง database

* core/port จะเก็บ interface ทั้งหมดของโปรเจ็คคอยกำหนดข้อมูลที่จะรับเเละค่าที่จะส่งกลับไปโดยมี port ทั้งหมด 2 port คือ inbound port เเละ outbound poort
    * inbound port : port ขาเข้ากำหนดการรับข้อมูลเข้ามาภายใน project ผ่าน api ว่าจะรับอะไรบ้างเพื่อส่งไปคำนวนที่ core/service ซึ่งเป็น business logic ของระบบโดยไฟล์ของ inbound port จะลงท้ายด้วย service เช่น core/port/user/user_service

    * outbound port : port ขาออกทำหน้าที่เป็น interface ของ database เพื่อให้ core service สามารถใช้งานกับ database ได้โดยไฟล์ของ outbound port จะลงท้ายด้วย repository เช่น core/port/user/user_repository

* adapter ทำหน้าที่ implement interface เพื่อนำไปใช้งานโดยเเบ่งตาม port เป็น inbound adapter เเละ outbound adapter 
    * inbound adapter : ทำหน้าที่จัดการ request ที่ถูกส่งเข้ามาเเปลงข้อมูลจาก json มาเป็น entities เเละส่งไปให้ core service 

    * outbound adapter : ทำหน้าที่จัดการกับข้อมูลที่ได้จาก core service มาเขียนลง database


* core/service ทำหน้าที่จัดการกับ business logic ทั้งหมดจะของโปรเเกรม

## How to use
    docker compose up -d

    go test ./... (optional)

    go run main.go

    * สามารถใช้งานได้เลย!!


## Api List
การเข้าถึง api จะถูกกำหนดสิทธิ์การเข้าถึงไว้โดยมี role กำกับทั้งหมด 3 role ดังนี้ 

> user :  สิทธิ์การใช้งานระดับต้นของระบบสามารถเพิ่ม ลบ เเก้ไข ข้อมูลของตัวเองเเละสามารถสั่งซื้อสินค้าได้ผู้ใช้ที่สร้างผ่าน api /register/user จะได้รับ role นี้โดยอัตโนมัต

> seller : สิทธิ์การใช้งานระดับกลางของระบบสามารถ เพิ่ม ลบ เเก้ไขสินค้าได้เเละข้อมูลของตัวเองได้เเละสั่งซื้อสินค้าได้ผู้ใช้ที่สร้างผ่าน api /register/seller จะได้รับ role นี้โดยอัตโนมัต

> admin : สิทธิ์การใช้งานระดับสูงสุดของระบบสามารถเข้าถึงได้ทุกการทำงานของระบบสามารถเพิ่ม ลบ เเก้ไข ผู้ใช้เเละสินค้าได้ทั้งหมด

* ### User (สำหรับจัดการกับ user)  
    * GET
        * ``` localhost:3030/users ```
        *เเสดงข้อมูลผู้ใช้ทั้งหมด Role : admin
        * ``` localhost:3030/users/:id ``` 
        *เเสดงข้อมูลผู้ใช้ตาม id Role : user, seller, admin
        
    * POST
        * ``` localhost:3030/register/user ``` 
        *สร้างผู้ใช้สิทธิ์ user
        * ``` localhost:3030/register/seller ``` 
        *สร้างผู้ใช้สิทธิ์ seller
        * ``` localhost:3030/login ``` 
        *เข้าสู่ระบบ

    * PATCH
        * ``` localhost:3030/users/:id ```
        *เเก้ไขข้อมูลผู้ใช้ Role : user, seller, admin

    * DELETE
        * ``` localhost:3030/users/:id ``` 
        *ลบผู้ใช้ Role : user, seller, admin

* ### Product (สำหรับจัดการกับ product)
    * GET
        * ``` localhost:3030/products ``` 
        *เเสดงข้อมูลสินค้าทั้งหมด Role : user, seller, admin
        * ``` localhost:3030/products/:id ``` 
        *เเสดงข้อมูลสินค้าตาม id Role : user, seller, admin
        
    * POST 
        * ``` localhost:3030/products ``` 
        *สร้างสินค้า Role : seller ,admin

    * PATCH
        * ``` localhost:3030/products/:id ```
        *เเก้ไขสินค้า Role : seller, admin

    * DELETE
        * ``` localhost:3030/products/:id ```
        *ลบสินค้า Role : seller ,admin

* ### Order (สำหรับจัดการกับ order)
    * GET
        * ``` localhost:3030/orders/user/:id ``` *เเสดงข้อมูลคำสั่งซื้อทั้งหมดของผู้ใช้งานตาม id ของผู้ใช้ Role : user, seller, admin
        
    * POST 
        * ``` localhost:3030/orders ```
        *สร้างคำสั่งซื้อ Role : user, seller ,admin

    * PATCH
        * ``` localhost:3030/orders/:id ``` 
        *เเก้ไขคำสั่งซื้อค้า Role : user, seller, admin

    * DELETE
        * ``` localhost:3030/orders/:id ```
        *ลบคำสั่งซื้อ Role : user, seller ,admin

## Software Require
* docker 
* go version > 1.22
