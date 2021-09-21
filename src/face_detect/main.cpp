#include<opencv2/opencv.hpp>
#include<filesystem>
#include<iostream>
#include<random>
#include "cv_algorithm.hpp"

constexpr int offset_width = 128;
constexpr int offset_height = 64;

//HSV?F????p?l
constexpr int B_MIN = 80;
constexpr int B_MAX = 200;
constexpr int G_MIN = 80;
constexpr int G_MAX = 220;
constexpr int R_MIN = 160;
constexpr int R_MAX = 255;
constexpr int H_MAX = 30;
constexpr int H_MIN = 0;
constexpr int S_MAX = 150;
constexpr int S_MIN = 30;
constexpr int V_MAX = 255;
constexpr int V_MIN = 60;

#define MIN_HSVCOLOR cv::Scalar(0, 60, 80)
#define MAX_HSVCOLOR cv::Scalar(10, 160, 240)

#define MIN_BKCOLOR cv::Scalar(0, 0, 0)
#define MAX_BKCOLOR cv::Scalar(180, 255, 40)

//1箱の限界値を定める
constexpr int img_height = 500;
constexpr int img_width = 500;
constexpr double resize_param = 5.2;
constexpr double resize_center = resize_param * 2;

cv::Mat PinP_point(const cv::UMat &srcImg, const cv::UMat &smallImg, const cv::Point2f p0, const cv::Point2f p1)
{

	cv::Mat dstImg;
	srcImg.copyTo(dstImg);

	std::vector<cv::Point2f> src, dst;
	src.push_back(cv::Point2f(0, 0));
	src.push_back(cv::Point2f(smallImg.cols, 0));
	src.push_back(cv::Point2f(smallImg.cols, smallImg.rows));

	dst.push_back(p0);
	dst.push_back(cv::Point2f(p1.x, p0.y));
	dst.push_back(p1);

	cv::Mat mat = cv::getAffineTransform(src, dst);

	cv::warpAffine(smallImg, dstImg, mat, dstImg.size(), cv::INTER_LINEAR, cv::BORDER_TRANSPARENT);
	return dstImg;
}

int main(int argc,char*argv[])
{
	std::string filePath;
	if(argc < 2){
		filePath = "/template_1.jpg";
	}else{
		filePath = argv[1];
	}
	std::string filename;
	cv::UMat result_img;

	cv::Mat template_img = cv::imread(filePath);
	template_img.copyTo(result_img);
	int offset_width = 0;
	int offset_height = 0;
	//cv::Mat test_img = cv::imread("images/test_img.png");
	int width = 100;
	int height = 100;
	int index_w = 0;
	int index_h = 0;
	int size_w = template_img.cols;
	int size_h = template_img.rows;

	std::random_device rnd;     // 非決定的な乱数生成器でシード生成機を生成
	std::mt19937 mt(rnd()); //  メルセンヌツイスターの32ビット版、引数は初期シード
	//ここで画像の配置候補を作成する

	int cnt_h = static_cast<int>(template_img.rows / img_height);
	int cnt_w = static_cast<int>(template_img.cols / img_width);

	std::uniform_int_distribution<>rand_w(0, cnt_w - 2);     // [0, 99] 範囲の一様乱数
	std::uniform_int_distribution<>rand_h(0, cnt_h - 2);     // [0, 99] 範囲の一様乱数
	std::vector<std::vector<bool>> img_map(cnt_w, std::vector<bool>(cnt_h));

	for (const std::filesystem::directory_entry &i : std::filesystem::directory_iterator("/cheese/images")){

		filename = "/cheese/images/" + i.path().filename().string();

		std::cout << i.path().filename().string() << std::endl;
		cv::UMat extract_img;

		extractFaceImage(filename, extract_img);
		if (extract_img.empty()){
			std::cout << "skipped" << std::endl;
			continue;
		}
		//配置場所を決める

		do {
			index_w = rand_w(rnd);
			index_h = rand_h(rnd);
		} while (img_map.at(index_w).at(index_h));

		img_map.at(index_w).at(index_h) = false;

		cv::circle(result_img, cv::Point2f(500.0 * index_w + (extract_img.cols / resize_center), 500.0 * index_h + (extract_img.rows / resize_center)), 300, cv::Scalar(240, 240, 240),-1);
		auto tmp = PinP_point(result_img, extract_img, cv::Point2f(500.0 * index_w, 500.0 * index_h), cv::Point2f(500.0 * index_w + (extract_img.cols / resize_param), 500.0 * index_h + (extract_img.rows / resize_param)));
		tmp.copyTo(result_img);
	}
	cv::imwrite("/cheese/result.jpg", result_img);
}
